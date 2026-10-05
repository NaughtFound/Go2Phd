package main

import (
	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/handlers"
	"naughtfound.github.io/go2phd/repositories"
	"naughtfound.github.io/go2phd/services"
)

func main() {
	r := gin.Default()
	userGroups := r.Group("users")

	repo := repositories.NewUserRepo()
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	handler.SetupGroup(userGroups)

	r.Run(":8080")
}
