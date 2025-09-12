package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_ai/domain"
	"io"
	"net/http"
	"path"
)

const serverUrl = "http://127.0.0.1:8080/"

var (
	availableModelsUrl = path.Join(serverUrl, "models")
)

func get[Resp any](url string) (Resp, error) {
	var resp Resp
	data, err := http.Get(url)
	if err != nil {
		return resp, err
	}
	defer data.Body.Close()

	body, err := io.ReadAll(data.Body)
	if err != nil {
		return resp, err
	}

	if data.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("сервер вернул ошибку %d: %s", data.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func post(model, prompt string) (domain.Response, error) {
	req := domain.Request{
		Model: model,
		Input: prompt,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return domain.Response{}, err
	}

	resp, err := http.Post(serverUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return domain.Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Response{}, fmt.Errorf("сервер вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	domainResp := domain.Response{}
	if err := json.Unmarshal(body, &domainResp); err != nil {
		return domain.Response{}, err
	}

	return domainResp, nil
}
