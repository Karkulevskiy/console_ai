package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func startServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", askAi)
	e.POST("/tries", askAiWithManyTries)
	e.GET("/models", getModelsHandler)
	e.POST("/models", addModelHandler)
	e.DELETE("/models", deleteModelHandler)
	e.PUT("/models", updateModelHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	<-quit
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server stopped.")
}
