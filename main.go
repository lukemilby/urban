package main

import (
	"fmt"
	"github.com/dpatrie/urbandictionary"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"sort"
	"strings"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	w, h := ui.TerminalDimensions()

	search := widgets.NewParagraph()
	search.Title = "Urban Search"
	search.Text = ""
	search.SetRect( 0, 0, w, 3)
	search.BorderStyle.Fg = ui.ColorYellow

	urbanList := widgets.NewList()
	urbanList.Title = "Results"
	urbanList.SetRect(0,3,w,h )
	urbanList.BorderStyle.Fg = ui.ColorYellow
	ui.Render(search, urbanList)
	uiEvents := ui.PollEvents()

	var searchResults []urbandictionary.Result

	for {
		w, h := ui.TerminalDimensions()
		e := <-uiEvents

		switch e.ID {
		case "<C-c>":
			return
		case "<C-v>":
			search.Text = ""
			ui.Render(search, urbanList)
		case "<MouseLeft>":
			ui.Render(search, urbanList)
		case "<MouseRelease>":
			ui.Render(search, urbanList)
		case "<Backspace>":
			if len(search.Text) > 0 {
				search.Text = search.Text[:len(search.Text)-1]
			}
			ui.Render(search, urbanList)
		case "<Resize>":
			urbanList.SetRect(0,3,w,h-1)
			ui.Render(search, urbanList)
		case "<Enter>":
			if len(urbanList.Rows) == 0 {
				searchResults = Lookup(search.Text)
				var payload []string
				for _, r := range searchResults {
					payload = append(payload, strings.Replace(
						fmt.Sprintf("Upvotes: %d Definition: %v ", r.Upvote, r.Definition), "\r\n", "", -1))
				}
				urbanList.Rows = payload
				urbanList.SelectedRowStyle = ui.Style{
					Fg:       0,
					Bg:       2,
					Modifier: 0,
				}
				ui.Render(search, urbanList)
			} else {
				selectedItem := searchResults[urbanList.SelectedRow]
				selectedI := widgets.NewParagraph()
				selectedI.Title = selectedItem.Word
				selectedI.Text = selectedItem.Definition
				selectedI.WrapText = true
				selectedI.SetRect( 10, 2, w-20, 10)
				selectedI.BorderStyle.Fg = ui.ColorYellow
				ui.Render(selectedI)
			}

		case "<Down>":
			urbanList.SelectedRowStyle = ui.Style{
				Fg:       0,
				Bg:       2,
				Modifier: 0,
			}
			urbanList.ScrollDown()
			ui.Render(search, urbanList)
		case "<Up>":
			urbanList.SelectedRowStyle = ui.Style{
				Fg:       0,
				Bg:       2,
				Modifier: 0,
			}
			urbanList.ScrollUp()
			ui.Render(search, urbanList)
		case "<Escape>":
			ui.Render(search, urbanList)
		default:
			search.Text += e.ID
			ui.Render(search, urbanList)
		}
	}
}

func Lookup(s string) []urbandictionary.Result {
	sr, err := urbandictionary.Query(s)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(sr.Results, func(i, j int) bool {
		return sr.Results[i].Upvote > sr.Results[j].Upvote
	})

	return sr.Results
}