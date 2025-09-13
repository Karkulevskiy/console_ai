package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const serverUrl = "http://127.0.0.1:8080/"

var (
	availableModelsUrl = serverUrl + "models"
	modelApiKeyUrl     = serverUrl + "model/api_key"
)

func getStr(url string) (string, error) {
	bytes, err := get(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getWithType[Resp any](url string) (Resp, error) {
	var resp Resp
	bytes, err := get(url)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func get(url string) ([]byte, error) {
	data, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	bytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	if data.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер вернул ошибку %d: %s", data.StatusCode, string(bytes))
	}
	return bytes, nil
}

func put[Request any](url string, req Request) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func post[Request, Response any](req Request, url string) (Response, error) {
	var resp Response

	data, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	postResp, err := http.Post(serverUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return resp, err
	}
	defer postResp.Body.Close()

	body, err := io.ReadAll(postResp.Body)
	if err != nil {
		return resp, err
	}

	if postResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("сервер вернул ошибку %d: %s", postResp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
