package main

import (
	"example/controller"
	"example/db"
	"example/repository"
	"example/service"
	"github.com/aivyss/bean"
)

func main() {
	buf := bean.GetBeanBuffer()

	buf.RegisterBean(controller.NewUserController)
	buf.RegisterBean(repository.NewUserRepository)
	buf.RegisterBean(service.NewUserService)
	buf.RegisterBean(db.NewDB)

	if err := buf.Buffer(); err != nil {
		panic(err)
	}
}
