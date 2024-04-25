package main

import (
	"fmt"
	"os"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/handler"
	"github.com/nvdaz/find-a-friend-api/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database, err := db.NewDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer database.Close()

	userStore := db.NewUserStore(database.Db)
	serviceConversationStore := db.NewServiceConversationStore(database.Db)
	userService := service.NewUserService(userStore, serviceConversationStore)
	h := handler.NewHandler(userService, &serviceConversationStore)

	e := echo.New()

	e.Use(middleware.CORS())
	e.POST("/user", h.CreateUser)
	e.GET("/user/:id", h.GetUser)
	e.POST("/user/:id", h.UpdateUser)
	e.GET("/user/:id/match", h.GetUserMatches)
	e.GET("/users", h.GetAllUsers)
	e.POST("/service-conversations", h.CreateServiceConversations)
	e.GET("/service-conversations/:id", h.GetServiceConversations)

	fmt.Println("Starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
