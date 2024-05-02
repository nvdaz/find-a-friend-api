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

	userStore := db.NewUserStore(database)
	messageStore := db.NewMessagesStore(database)
	matchStore := db.NewMatchStore(database)
	userService := service.NewUserService(userStore, messageStore)
	matchService := service.NewMatchService(userService, matchStore)
	messageService := service.NewMessagesService(messageStore, userService)
	h := handler.NewHandler(userService, matchService, messageService)

	e := echo.New()

	e.Use(middleware.CORS())
	e.POST("/user", h.CreateUser)
	e.GET("/user/:id", h.GetUser)
	e.POST("/user/:id", h.UpdateUser)
	e.GET("/user/:id/matches", h.GetUserMatches)
	e.POST("/user/:id/matches", h.GenerateUserMatch)
	e.GET("/users", h.GetAllUsers)
	e.GET("/match/:id", h.GetMatch)
	e.POST("/messages", h.GetMessages)
	e.POST("/messages/create", h.CreateMessage)
	e.POST("/messages/poll", h.PollMessages)
	e.POST("/register", h.RegisterUser)
	e.POST("/login", h.LoginUser)
	e.POST("/set-icon", h.UpdateUserIcon)

	fmt.Println("Starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
