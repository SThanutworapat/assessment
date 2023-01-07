package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SThanutworapat/assessment/config"
	"github.com/SThanutworapat/assessment/db"
	expenses "github.com/SThanutworapat/assessment/handlers"
	expense "github.com/SThanutworapat/assessment/models"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	d := db.InitDB()
	e := echo.New()
	h := expenses.NewHandler(expense.NewDatabaseModel(d))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.POST("/expenses", h.CreateExpensesHandler)
	e.GET("/expenses", h.GetExpensesHandler)
	e.PUT("/expenses/:id", h.PutExpensesByIdHandler)
	e.GET("/expenses/:id", h.GetExpensesByIdHandler)
	e.POST("/quit", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		defer func() {
			log.Println("-----shutting down the server--------")
		}()
		return e.Shutdown(ctx)
	})

	// Start server
	config := config.NewConfig()
	go func() {
		if err := e.Start(":" + config.Port); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("-----shutting down the server--------")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		log.Println("-----shutting down the server--------")
	}()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
