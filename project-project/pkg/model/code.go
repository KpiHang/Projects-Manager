package model

import (
	"test.com/project-common/errs"
)

//const (
//	NoLegalMobile common.BusinessCode = 2001 // 手机号不合法代码
//)

var (
	RedisError  = errs.NewError(999, "redis 错误")
	DBError     = errs.NewError(998, "DB 错误")
	ParamsError = errs.NewError(401, "参数错误")

	NoLegalMobile         = errs.NewError(10102001, "手机号不合法") // 10 user模块 10 登录相关
	CaptchaNotExist       = errs.NewError(10102002, "验证码不存在，或者已过期")
	CaptchaError          = errs.NewError(10102003, "验证码错误")
	EmailExist            = errs.NewError(10102004, "邮箱已经存在")
	AccountExist          = errs.NewError(10102005, "账号已经存在")
	MobileExist           = errs.NewError(10102006, "手机号已经存在")
	AccountOrPwdError     = errs.NewError(10102007, "账号密码不正确")
	TaskNameNotNull       = errs.NewError(20102001, "任务标题不能为空")
	TaskStagesNotNull     = errs.NewError(20102002, "任务步骤不能为空")
	ProjectAlreadyDeleted = errs.NewError(20102003, "项目已被删除")
)
