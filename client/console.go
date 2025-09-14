package main

import (
	"context"
	"fmt"
	"go_ai/domain"
	"go_ai/logging"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var currApiKey = ""

func newModelList() *tview.List {
	models, err := getAvailableModels()
	if err != nil {
		panic(err)
	}
	modelList := tview.NewList().ShowSecondaryText(false)
	modelList.SetBorder(true).SetTitle("Выбор модели")
	for _, m := range models {
		modelList.AddItem(m, "", 0, nil)
	}
	modelList.SetCurrentItem(0)
	return modelList
}

func newResponseBox() *tview.TextView {
	responseBox := tview.NewTextView().SetDynamicColors(true)
	responseBox.SetBorder(true).SetTitle("Ответ").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlY {
				text := responseBox.GetText(true)
				if err := clipboard.WriteAll(text); err != nil {
					logging.Log(err.Error())
				}
				return nil
			}
			return event
		})
	return responseBox
}

func newInput() *tview.TextArea {
	input := tview.NewTextArea().SetPlaceholder("Введите промпт...")
	input.SetBorder(true).SetTitle("Промт:")
	return input
}

func newHelpBar() *tview.TextView {
	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]TAB/Shift+TAB[-]: переключение фокуса   " +
			"[yellow]Ctrl+Y[-]: копировать ответ   " +
			"[yellow]ESC[-]: выход")
	help.SetBorder(true).SetTitle("Подсказки")
	return help
}

func setPromptInputCapture(app *tview.Application, flex *tview.Flex, input *tview.TextArea, modelList *tview.List, responseBox *tview.TextView) {
	inHandler := newInputHandler()
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyEnter {
			return event
		}
		prompt := input.GetText()
		model, _ := modelList.GetItemText(modelList.GetCurrentItem())
		if strings.TrimSpace(prompt) == "" {
			return nil
		}
		if strings.HasPrefix(model, "googleai") {
			// TODO это работает пока что только для google херни
			callback := func() {
				app.SetFocus(flex)
			}
			ensureModelAPIKey(app, flex, model, callback)
		}
		out := inHandler(context.Background(), prompt, model, currApiKey)
		input.SetText("", true)
		responseBox.SetText(out)
		return nil
	})
}

func setAppInputCapture(app *tview.Application, focusables []tview.Primitive) {
	focusIndex := 0
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			focusIndex = (focusIndex + 1) % len(focusables)
			app.SetFocus(focusables[focusIndex])
			return nil
		case tcell.KeyBacktab:
			focusIndex = (focusIndex - 1 + len(focusables)) % len(focusables)
			app.SetFocus(focusables[focusIndex])
			return nil
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})
}

func getModelUrl(modelName string) string {
	params := url.Values{"model": []string{modelName}}
	params.Encode()
	return modelApiKeyUrl + "?" + params.Encode()
}

func ensureModelAPIKey(app *tview.Application, flex *tview.Flex, modelName string, callback func()) {
	if currApiKey != "" {
		return
	}
	apiKey, err := getStr(getModelUrl(modelName))
	if err != nil {
		logging.Log(fmt.Sprintf("failed to get api key: %v", err.Error()))
	}
	if apiKey = strings.TrimSpace(apiKey); apiKey != "" && len(apiKey) > 5 {
		currApiKey = apiKey
		if err := os.Setenv("GEMINI_API_KEY", currApiKey); err != nil {
			logging.Log(fmt.Sprintf("failed to set env var: %v", err.Error()))
		}
		callback()
		return
	}

	form := tview.NewForm().AddPasswordField("Введите API ключ для "+modelName, "", 40, '*', nil)
	form.AddButton("OK", func() {
		defer callback()
		inputApiKey := form.GetFormItem(0).(*tview.InputField).GetText()
		inputApiKey = strings.TrimSpace(inputApiKey)

		if inputApiKey == "" {
			return
		}

		if err := put(modelApiKeyUrl, domain.Model{Model: modelName, APIKey: inputApiKey}); err != nil {
			logging.Log(err.Error())
		}

		app.SetRoot(flex, true).SetFocus(flex)
		currApiKey = inputApiKey
		if err := os.Setenv("GEMINI_API_KEY", currApiKey); err != nil {
			logging.Log(fmt.Sprintf("failed to set env var: %v", err.Error()))
		}
	})

	form.SetBorder(true).SetTitle("API Key Required").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).SetFocus(form)
}
