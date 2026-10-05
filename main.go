package main

import (
	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/handlers"
	"naughtfound.github.io/go2phd/repositories"
	"naughtfound.github.io/go2phd/services"
)

func main() {
	r := gin.Default()
	universityGroups := r.Group("universities")

	repo := repositories.NewUniversityRepository()
	service := services.NewUniversityService(repo)
	handler := handlers.NewUniversityHandler(service)

	handler.SetupGroup(universityGroups)

	r.Run(":8080")
}
