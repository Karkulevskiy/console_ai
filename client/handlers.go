package main

import (
	"bytes"
	"context"
	"fmt"
	"go_ai/domain"
	"go_ai/logging"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func newInputHandler() func(context.Context, string, string, string) string {
	return func(ctx context.Context, prompt, model, apiKey string) string {
		req := domain.Request{
			Input:  prompt,
			Model:  model,
			APIKey: apiKey,
		}
		out, err := post[domain.Request, domain.Response](req, serverUrl)
		if err != nil {
			return "oops, smth went wrong :("
		}

		if strings.TrimSpace(out.Code) != "" {
			highlightedCode, err := highlightCode(out)
			if err != nil {
				logging.Log(err.Error())
			} else {
				out.Code = addTabsToCode(highlightedCode)
			}
		}

		answer := fmt.Sprintf("[yellow]Модель:[white] %s\n[yellow]Промпт:[white] %s\n[yellow]Ответ:[white]\n{%s\n}",
			model, prompt, out.String())

		return answer
	}
}

func getAvailableModels() ([]string, error) {
	return getWithType[[]string](availableModelsUrl)
}

func addTabsToCode(code string) string {
	var sb strings.Builder
	splittedCode := strings.SplitSeq(code, "\n")
	for line := range splittedCode {
		sb.WriteString("\t\t" + line + "\n")
	}
	return sb.String()
}

func highlightCode(r domain.Response) (string, error) {
	lexer := lexers.Get(r.Lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, r.Code)
	if err != nil {
		return "", err
	}

	style := styles.Get(styles.DoomOne2.Name)
	if style == nil {
		style = styles.Fallback
	}

	var buf bytes.Buffer
	for token := iterator(); token != chroma.EOF; token = iterator() {
		entry := style.Get(token.Type)
		color := entry.Colour.String()
		if color != "" {
			buf.WriteString(fmt.Sprintf("[%s]", color))
		}
		buf.WriteString(token.Value)
		if color != "" {
			buf.WriteString("[-]")
		}
	}

	return buf.String(), nil
}
