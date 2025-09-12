package main

import (
	"fmt"
	"go_ai/ai"
	"go_ai/db"
	"go_ai/domain"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Тогда должен быть еще какой то клиент, который делает POST запросы сюды через cli

func mockResponse() (domain.Response, error) {
	return domain.Response{Output: "Это замоканный респонс, азазазаз"}, nil
}

func askAi(c echo.Context) error {
	// Придумать стандратный формат инпута
	// Его парсинг, валидация
	// Запрос к гемини

	req := domain.Request{}

	if err := c.Bind(&req); err != nil {
		slog.Error("failed to unmarshall request")
		return err
	}

	fmt.Println("Input")
	fmt.Println(req.Input)
	resp, err := mockResponse()
	if err != nil {

	}
	return c.JSON(http.StatusOK, resp)

	// resp, err := ai.AskAI(c.Request().Context(), req)
	// fmt.Println("RESPONSE:")
	// fmt.Println(resp)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	//
	// return c.JSON(http.StatusOK, resp)
}

func askAiWithManyTries(c echo.Context) error {
	req := domain.Request{}

	if err := c.Bind(&req); err != nil {
		slog.Error("failed to unmarshall request")
		return err
	}

	fmt.Println("Input")
	fmt.Println(req.Input)

	resp, err := ai.AskAIWithManyTries(c.Request().Context(), req)
	fmt.Println("RESPONSE:")
	fmt.Println(resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}

type ModelRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name,omitempty"`
}

func getModelsHandler(c echo.Context) error {
	models, err := db.GetAvailableModels(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, models)
}

func addModelHandler(c echo.Context) error {
	req := ModelRequest{}
	if err := c.Bind(&req); err != nil || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.AddModel(c.Request().Context(), req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "added"})
}

func deleteModelHandler(c echo.Context) error {
	req := ModelRequest{}
	if err := c.Bind(&req); err != nil || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.DeleteModel(c.Request().Context(), req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func updateModelHandler(c echo.Context) error {
	req := ModelRequest{}
	if err := c.Bind(&req); err != nil || req.Name == "" || req.NewName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.UpdateModel(c.Request().Context(), req.Name, req.NewName); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}
