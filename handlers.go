package main

import (
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
	req := domain.Request{}
	if err := c.Bind(&req); err != nil {
		slog.Error("failed to unmarshall request")
		return err
	}
	resp, err := ai.AskAI(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func askAiWithManyTries(c echo.Context) error {
	req := domain.Request{}
	if err := c.Bind(&req); err != nil {
		slog.Error("failed to unmarshall request")
		return err
	}
	resp, err := ai.AskAIWithManyTries(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func getModelsHandler(c echo.Context) error {
	models, err := db.GetAvailableModels(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, models)
}

func getModelHandler(c echo.Context) error {
	modelName := c.QueryParam("model")
	model, err := db.GetModel(c.Request().Context(), modelName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, model)
}

func addModelHandler(c echo.Context) error {
	req := domain.Request{}
	if err := c.Bind(&req); err != nil || req.Model == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.AddModel(c.Request().Context(), req.Model); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "added"})
}

func deleteModelHandler(c echo.Context) error {
	req := domain.Request{}
	if err := c.Bind(&req); err != nil || req.Model == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.DeleteModel(c.Request().Context(), req.Model); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func updateModelHandler(c echo.Context) error {
	return nil
	// req := domain.Request{}
	// if err := c.Bind(&req); err != nil || req.Name == "" || req.NewName == "" {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	// }
	// if err := db.UpdateModelName(c.Request().Context(), req.Name, req.NewName); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// }
	// return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func getModelAPIKeyHandler(c echo.Context) error {
	modelName := c.QueryParam("model")
	apiKey, err := db.GetModelAPIKey(c.Request().Context(), modelName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apiKey)
}

func updateModelAPIKeyHandler(c echo.Context) error {
	req := domain.Model{}
	if err := c.Bind(&req); err != nil || req.Model == "" || req.APIKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := db.UpdateModelAPIKey(c.Request().Context(), req.Model, req.APIKey); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "api_key updated"})
}
