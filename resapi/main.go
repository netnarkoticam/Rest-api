package main

import (
	"fmt"
	"os"
	"resapi/internal/config"
	"resapi/internal/repository"
	"resapi/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.New(db)
	srv := server.Server{Db: repo}

	router := gin.Default()

	router.GET("/users", srv.GetUsersHandler)
	router.POST("/users", srv.SaveUserHandler)
	router.PATCH("/users/:id", srv.PatchUserHandler)
	router.DELETE("/users/:id", srv.DeleteUserHandler)
	router.PUT("/users/:id", srv.UpdateUserHandler)

	// Читаем адрес из конфигурации (если нужно)
	// config := config.ReadConfig()
	// router.Run(config.Address)

	router.Run(":8080")
}