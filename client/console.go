package main

import (
	"context"
	"go_ai/logging"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

func setPromptInputCapture(input *tview.TextArea, modelList *tview.List, responseBox *tview.TextView) {
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
		out := inHandler(context.Background(), prompt, model)
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
