package main

import "github.com/rivo/tview"

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
	responseBox := tview.NewTextView().
		SetDynamicColors(true)
	responseBox.SetBorder(true).
		SetTitle("Ответ")
	return responseBox
}

func newInput() *tview.TextArea {
	input := tview.NewTextArea().SetPlaceholder("Введите промпт...")
	input.SetBorder(true).SetTitle("Промт:")
	return input
}
