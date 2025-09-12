package db

import "context"

func GetAvailableModels(ctx context.Context) ([]string, error) {
	// TODO
	return []string{
		"GPT-3.5",
		"GPT-4",
		"LLaMA-2",
		"Custom-AI",
		"googleai/gemini-2.5-flash",
	}, nil
}
