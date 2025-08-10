package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"net/http"
	//"fmt"
	"io"
	"log"
	"strings"
	"os"
)

func main() {
	app := tview.NewApplication()
	
	dataView := tview.NewTextView().
					SetDynamicColors(true).
					SetScrollable(true).
					SetRegions(true).
					SetWordWrap(true).
					SetChangedFunc(func() {
						app.Draw()
					})

	messageBox := tview.NewTextView().
						SetRegions(true).
						SetChangedFunc(func() {
							app.Draw()
						})

	targetField := tview.NewInputField().
					SetLabel("Target:").
					SetFieldWidth(10).
					SetFieldBackgroundColor(tcell.ColorRed)

	payloadField := tview.NewInputField().
					SetLabel("Payload:").
					SetFieldWidth(9).
					SetFieldBackgroundColor(tcell.ColorRed)

	sendButton := tview.NewButton("Send request").
			 			SetSelectedFunc(func() {
			 				craftedReq := targetField.GetText() + "/" + payloadField.GetText()

			 				if targetField.GetText() == "" {
			 					messageBox.SetText("[!] Please, input target url")
			 				} else {

								resp, err := http.Get(craftedReq)
								if err != nil {
									log.Fatal(err)
								}	
								defer resp.Body.Close()

								data, err := io.ReadAll(resp.Body)
								if err != nil {
									messageBox.SetText("[!] Something went wrong..")
									log.Fatal(err)
								}
				

								dataView.SetText(string(data))

								if strings.Contains(string(data), payloadField.GetText()) {
									messageBox.SetText("[$] It's seems that is page reflecting your payload!")
									strings.Replace(string(data), payloadField.GetText(), 
																"[#ff0000]"+payloadField.GetText()+"[white]", -1)
								}
							}
						})

	clearButton := tview.NewButton("Clear responce").
						SetSelectedFunc(func() {
							dataView.SetText(" ")
							messageBox.SetText("[#] There is nothing yet...")
						})

	saveButton := tview.NewButton("Save responce").
						SetSelectedFunc(func() {
							file, err := os.OpenFile("responce.html", os.O_WRONLY|os.O_CREATE, 0644)
							if err != nil {
								log.Fatal(err)
							}

							_, err = file.WriteString(dataView.GetText(true))
							if err != nil {
								log.Fatal(err)
							}

							messageBox.SetText("[#] Responce saved")
						})

	inputPanel := tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(targetField, 0, 1, true).
						AddItem(payloadField, 0, 1, true)

	buttonPanel := tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(sendButton, 0, 1, false).
						AddItem(clearButton, 0, 1, false)

	saveFlex := tview.NewFlex().
					AddItem(saveButton, 0, 1, false)

	informationPanes := tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(dataView, 0, 12, true).
							AddItem(saveFlex, 0, 3, false).
							AddItem(messageBox, 0, 2, false)
						
		
	secondFlex := tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(inputPanel, 0, 2, false).
						AddItem(buttonPanel, 0, 5, false)
	
	mainFlex := tview.NewFlex().
						AddItem(informationPanes, 0, 1, false).
						AddItem(secondFlex, 20, 1, true)

	inputPanel.SetBorder(true)
	inputPanel.SetTitle("Input")

	buttonPanel.SetBorder(true)
	sendButton.SetBorder(true)
	saveButton.SetBorder(true)
	saveFlex.SetBorder(true)
	clearButton.SetBackgroundColorActivated(tcell.ColorRed).SetBorder(true)
	
	dataView.SetBorder(true)
	dataView.SetTitle("Responce")
	dataView.SetText(" ")

	messageBox.SetBorder(true)
	messageBox.SetTitle("Info messages")
	messageBox.SetText("[#] There is nothing yet...")


	err := app.SetRoot(mainFlex, true).EnableMouse(true).Run()
	if err != nil {
		messageBox.SetText("[!] Something went wrong..")
		log.Fatal(err)
	}
	
}