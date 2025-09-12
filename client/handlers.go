package main

import (
	"context"
	"fmt"
)

func newInputHandler() func(context.Context, string, string) string {
	return func(ctx context.Context, prompt, model string) string {
		out, err := post(model, prompt)
		if err != nil {
			return "oops, smth went wrong :("
		}

		answer := fmt.Sprintf("[yellow]Модель:[white] %s\n[yellow]Промпт:[white] %s\n[yellow]Ответ:[white] %s",
			model, prompt, out.String())

		return answer
	}
}

func getAvailableModels() ([]string, error) {
	return get[[]string](availableModelsUrl)
}
