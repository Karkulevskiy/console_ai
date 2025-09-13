package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	modelList := newModelList()
	responseBox := newResponseBox()
	input := newInput()
	helpBar := newHelpBar()

	setPromptInputCapture(input, modelList, responseBox)
	setAppInputCapture(app, []tview.Primitive{modelList, responseBox, input})

	flex := tview.NewFlex().
		AddItem(modelList, 20, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(responseBox, 0, 3, false).
			AddItem(input, 0, 1, false).
			AddItem(helpBar, 3, 1, false),
			0, 3, false)

	// TODO
	// Добавить изменение цвета, что ожидаем ответ
	// Добавить экспорт api ключика, ну и собственно проверку этого
	if err := app.SetRoot(flex, true).
		EnableMouse(false).
		EnablePaste(true).
		Run(); err != nil {
		panic(err)
	}
}
