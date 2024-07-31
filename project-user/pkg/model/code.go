package model

import (
	"test.com/project-common/errs"
)

//const (
//	NoLegalMobile common.BusinessCode = 2001 // 手机号不合法代码
//)

var (
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)
