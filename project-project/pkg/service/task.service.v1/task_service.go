package task_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/tms"
	"test.com/project-grpc/task"
	"test.com/project-grpc/user/login"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/pkg/model"
	"time"
)

type TaskService struct {
	task.UnimplementedTaskServiceServer                  // grpc 里的
	cache                               repo.Cache       // repo里定义接口
	transaction                         tran.Transaction // 事务操作接口
	projectRepo                         repo.ProjectRepo
	projectTemplateRepo                 repo.ProjectTemplateRepo
	taskStagesTemplateRepo              repo.TaskStagesTemplateRepo
	taskStagesRepo                      repo.TaskStagesRepo
	taskRepo                            repo.TaskRepo
}

func NewTaskService() *TaskService {
	return &TaskService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransactionImpl(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		taskRepo:               dao.NewTaskDao(),
	}
}

func (t *TaskService) TaskStages(co context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	page := msg.Page
	pageSize := msg.PageSize

	// 根据projectCode查出所有的taskstages，在ms_task_stages表里
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stages, total, err := t.taskStagesRepo.FindStagesByProjectId(ctx, projectCode, page, pageSize)
	if err != nil {
		zap.L().Error("project TaskStages taskStagesRepo.FindStagesByProjectId DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}

	var tsMessages []*task.TaskStagesMessage
	copier.Copy(&tsMessages, stages)
	if tsMessages == nil {
		return &task.TaskStagesResponse{List: tsMessages, Total: 0}, nil
	}
	stagesMap := data.ToTaskStagesMap(stages)
	for _, v := range tsMessages {
		taskStage := stagesMap[int(v.Id)]
		v.Code = encrypts.EncryptNoErr(int64(v.Id))
		v.CreateTime = tms.FormatByMill(taskStage.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &task.TaskStagesResponse{List: tsMessages, Total: total}, nil
}

func (t *TaskService) MemberProjectList(co context.Context, msg *task.TaskReqMessage) (*task.MemberProjectResponse, error) {
	// 1. 去project_member表中去查找用户id列表；
	// 2. 用用户idlist 去查[]用户信息；

	// 1. 去project_member表中去查找用户id列表；
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	projectMembers, total, err := t.projectRepo.FindProjectMemberByPId(ctx, projectCode)
	if err != nil {
		zap.L().Error("project MemberProjectList projectRepo.FindProjectMemberByPId DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	// 2. 用用户idlist 去查[]用户信息；
	if projectMembers == nil || len(projectMembers) == 0 {
		return &task.MemberProjectResponse{List: nil, Total: 0}, nil
	}
	var mIds []int64
	pmMap := make(map[int64]*pro.ProjectMember)
	for _, v := range projectMembers {
		mIds = append(mIds, v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	// 请求用户信息；
	userMsg := &login.UserMessage{
		MIds: mIds,
	}
	memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, userMsg)
	if err != nil {
		zap.L().Error("project MemberProjectList LoginServiceClient.FindMemInfoByIds DB error, ", zap.Error(err)) // 非业务错误；
		return nil, err
	}
	var list []*task.MemberProjectMessage
	for _, v := range memberMessageList.List {
		owner := pmMap[v.Id].IsOwner
		mpm := &task.MemberProjectMessage{
			MemberCode: v.Id,
			Name:       v.Name,
			Avatar:     v.Avatar,
			Email:      v.Email,
			Code:       v.Code,
		}
		if v.Id == owner { // 默认值是0
			mpm.IsOwner = 1
		}
		list = append(list, mpm)
	}
	return &task.MemberProjectResponse{List: list, Total: total}, nil
}

func (t *TaskService) TaskList(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskList, err := t.taskRepo.FindTaskByStageCode(c, int(stageCode))
	if err != nil {
		zap.L().Error("project task TaskList FindTaskByStageCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var taskDisplayList []*data.TaskDisplay // task 转换为 taskDisplay 有的字段数据类型需要变化；
	var mIds []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			// 代表开启了隐私模式，不是每个用户都能见到这个task；
			taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
				return nil, errs.GrpcError(model.DBError)
			}
			if taskMember != nil {
				display.CanRead = model.CanRead
			} else {
				display.CanRead = model.NoCanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mIds = append(mIds, v.AssignTo)
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		zap.L().Error("project task TaskList rpc.LoginServiceClient.FindMemInfoByIds", zap.Error(err))
		return nil, err
	}
	if mIds == nil || len(mIds) <= 0 {
		return &task.TaskListResponse{List: nil}, nil
	}
	memberMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		memberMap[v.Id] = v
	}
	for _, v := range taskDisplayList {
		message := memberMap[encrypts.DecryptNoErr(v.AssignTo)]
		e := data.Executor{
			Name:   message.Name,
			Avatar: message.Avatar,
		}
		v.Executor = e
	}

	var taskMessageList []*task.TaskMessage
	copier.Copy(&taskMessageList, taskDisplayList)
	return &task.TaskListResponse{List: taskMessageList}, nil
}

func (t *TaskService) SaveTask(ctx context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	// 1. 检查业务逻辑；
	if msg.Name == "" {
		return nil, errs.GrpcError(model.TaskNameNotNull) // 任务名字不能为空
	}
	stageCode := encrypts.DecryptNoErr(msg.StageCode)
	taskStages, err := t.taskStagesRepo.FindById(ctx, int(stageCode)) // 确定stage存在；
	if err != nil {
		zap.L().Error("project task SaveTask taskStagesRepo.FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if taskStages == nil {
		return nil, errs.GrpcError(model.TaskStagesNotNull)
	}
	projectCode := encrypts.DecryptNoErr(msg.ProjectCode)
	project, err := t.projectRepo.FindProjectById(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask projectRepo.FindProjectById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if project == nil || project.Deleted == model.Deleted {
		return nil, errs.GrpcError(model.ProjectAlreadyDeleted)
	}
	// 2. 检查完了，可以保存任务task；

	// 设置字段maxIdNum、maxSort
	maxIdNum, err := t.taskRepo.FindTaskMaxIdNum(ctx, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskMaxIdNum error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxIdNum == nil {
		a := 0
		maxIdNum = &a
	}
	maxSort, err := t.taskRepo.FindTaskSort(ctx, projectCode, stageCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskSort error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	// 设置字段maxIdNum、maxSort

	// 保存
	assignTo := encrypts.DecryptNoErr(msg.AssignTo)
	ts := &data.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1, // 任务的id，递增的效果，每次最大的+1
		Private:     project.OpenTaskPrivate,
		Sort:        *maxSort + 1,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	// 保存操作要用事务
	err = t.transaction.Action(func(conn database.DbConn) error {
		err = t.taskRepo.SaveTask(ctx, conn, ts)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTask error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		tm := &data.TaskMember{
			MemberCode: assignTo,
			TaskCode:   ts.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    model.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = model.Executor
		}
		err = t.taskRepo.SaveTaskMember(ctx, conn, tm)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTaskMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	display := ts.ToTaskDisplay()
	member, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.UserMessage{MemId: assignTo})
	if err != nil {
		return nil, err
	}
	display.Executor = data.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}
	tm := &task.TaskMessage{}
	copier.Copy(tm, display)
	return tm, nil
}
