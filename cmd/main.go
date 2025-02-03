package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-conference/internal/database"
	"go-conference/internal/handlers"
	"log"
	"path/filepath"
)

func main() {
	database.Connect() // Connects to the Postgres db
	database.Migrate()

	r := gin.Default()

	files, err := filepath.Glob("web/templates/*")
	if err != nil {
		log.Fatalf("Error finding templates: %v", err)
	}
	fmt.Println("Templates found:", files)

	r.LoadHTMLGlob("web/templates/*.html")
	r.Static("/static", "./web/static")

	r.GET("/", handlers.Home)
	r.GET("/tickets", handlers.TicketForm)
	r.POST("/buy", handlers.BuyTicket)

	r.Run(":8080")
}
