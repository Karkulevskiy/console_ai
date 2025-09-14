package ai

import (
	"context"
	"fmt"
	"go_ai/domain"
	"go_ai/logging"
	"log/slog"
	"os"
	"sync"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

// TODO функционал выбора модели!!
// Нужно выдавать только google модельки

var (
	modelToApis = map[string]string{}
	once        sync.Once
)

// TODO подумать про шифрование апи ключиков
func setEnvApiKey(modelName, newApiKey string) {
	oldApiKey, ok := modelToApis[modelName]
	if ok && oldApiKey != "" && newApiKey == oldApiKey {
		return
	}
	// TODO обновление apiKey у модели в БД
	// TODO пока не работает обновление
	once.Do(func() {
		if err := os.Setenv("GEMINI_API_KEY", newApiKey); err != nil {
			logging.Log(fmt.Sprintf("failed to set env var: %v", err.Error()))
		}
		modelToApis[modelName] = newApiKey
	})
}

func AskAI(ctx context.Context, req domain.Request) (domain.Response, error) {
	// return domain.Response{Output: "Mock output", Code: "print('Hello')"}, nil
	g := newGenkit(ctx, req)
	flow := newFlow(g)
	setEnvApiKey(req.Model, req.APIKey)
	resp, err := flow.Run(ctx, req)
	if err != nil {
		slog.Error("failed to run prompt")
		return domain.Response{}, err
	}
	return resp, nil
}

func AskAIWithManyTries(ctx context.Context, req domain.Request) (domain.Response, error) {
	return askAIWithManyTries(ctx, req, 10)
}

func newFlow(g *genkit.Genkit) *core.Flow[domain.Request, domain.Response, struct{}] {
	return genkit.DefineFlow(g, "default flow", func(ctx context.Context, req domain.Request) (domain.Response, error) {
		prompt := req.Input
		aiOutput, _, err := genkit.GenerateData[domain.Response](ctx, g, ai.WithPrompt(prompt), ai.WithMaxTurns(1))
		if err != nil {
			slog.Error(fmt.Sprintf("failed to generate recipe: %v", err))
		}
		return *aiOutput, nil
	})
}

func newGenkit(ctx context.Context, req domain.Request) *genkit.Genkit {
	return genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel(req.GetModel()),
	)
}

func askAIWithManyTries(ctx context.Context, req domain.Request, triesCount int) (domain.Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan domain.Response, triesCount)
	var wg sync.WaitGroup
	wg.Add(triesCount)

	for i := range triesCount {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					g := newGenkit(ctx, req)
					flow := newFlow(g)
					resp, err := flow.Run(ctx, req)
					if err != nil {
						slog.Error(fmt.Sprintf("failed to run flow: from G: %d, err: %v", i, err))
						return
					}
					select {
					case <-ctx.Done():
						return
					case ch <- resp:
						cancel()
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case resp, ok := <-ch:
		if !ok {
			return domain.Response{}, fmt.Errorf("all attempts failed")
		}
		return resp, nil
	case <-ctx.Done():
		return domain.Response{}, ctx.Err()
	}
}
