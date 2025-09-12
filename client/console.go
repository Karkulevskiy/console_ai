package main

import (
	"go_ai/logging"

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
