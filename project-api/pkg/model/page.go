package model

import "github.com/gin-gonic/gin"

type Page struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"pageSize" form:"pageSize"`
}

func (p *Page) Bind(c *gin.Context) {
	_ = c.ShouldBind(&p)
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 { // 没有传入时，默认一页10个；
		p.PageSize = 10
	}
}
