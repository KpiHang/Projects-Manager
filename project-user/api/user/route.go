package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-user/router"
)

func init() {
	log.Println("init user router")
	router.AddRouter(&RouterUser{})
}

type RouterUser struct {
}

func (*RouterUser) Register(r *gin.Engine) {
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
