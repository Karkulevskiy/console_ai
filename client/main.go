package main

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	modelList := newModelList()
	responseBox := newResponseBox()
	input := newInput()
	inHandler := newInputHandler()

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			prompt := input.GetText()
			model, _ := modelList.GetItemText(modelList.GetCurrentItem())
			if prompt != "" {
				out := inHandler(context.Background(), prompt, model)
				input.SetText("", true)
				responseBox.SetText(out)
			}
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(modelList, 20, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(responseBox, 0, 3, false).
			AddItem(input, 0, 1, false),
			0, 3, false)

	focusables := []tview.Primitive{modelList, responseBox, input}
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

	if err := app.SetRoot(flex, true).
		EnableMouse(false).
		EnablePaste(true).
		Run(); err != nil {
		panic(err)
	}
}
