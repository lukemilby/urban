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
	search.Title = "Search"
	search.Text = ""
	search.SetRect( 0, 0, w, 3)
	search.BorderStyle.Fg = ui.ColorYellow

	urbanList := widgets.NewList()
	urbanList.Title = "Results"
	urbanList.SetRect(0,3,w,h )
	urbanList.BorderStyle.Fg = ui.ColorYellow
	ui.Render(search, urbanList)

	uiEvents := ui.PollEvents()
	for {
		w, h := ui.TerminalDimensions()
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Backspace>":
			if len(search.Text) > 0 {
				search.Text = search.Text[:len(search.Text)-1]
			}
			ui.Render(search, urbanList)
		case "<Resize>":
			urbanList.SetRect(0,3,w,h-1)
			ui.Render(search, urbanList)
		case "<MountLeft>":
			return
		case "<MouseRelease>":
			return
		case "<Enter>":
			sr := Lookup(search.Text)
			var payload []string
			for _, r := range sr {
				fmt.Println(r.Definition)
				payload = append(payload,strings.Replace(fmt.Sprintf("Upvotes: %d Definition: %s", r.Upvote, r.Definition), "", "", -1))
			}
			urbanList.Rows = payload
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