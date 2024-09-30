package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TUIChoice string

type GopherModel struct {
	choices  TUIChoice
	cursor   int
	selected map[int]struct{}
	tuiInput textinput.Model
}

func NewGopherModel(tuiInput textinput.Model) *GopherModel {
	return &GopherModel{
		tuiInput: tuiInput,
		selected: make(map[int]struct{}),
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Check if it is a key press
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit

		case "up", "k":
			if a.GopherModel.cursor > 0 {
				a.GopherModel.cursor--
			}

		case "down", "j":
			if a.GopherModel.cursor < len(a.GopherModel.choices)-1 {
				a.GopherModel.cursor++

			}
		case "enter", " ":
			_, ok := a.GopherModel.selected[a.GopherModel.cursor]
			if ok {
				delete(a.GopherModel.selected, a.GopherModel.cursor)
			} else {
				a.GopherModel.selected[a.GopherModel.cursor] = struct{}{}
			}
		}

	}
	return a, nil
}

func (a *App) View() string {
	s := "What table are you working on? \n\n"
	choices := a.Tables

	choiceList := []string{}
	for _, v := range choices {
		choiceList = append(choiceList, v.GetTableName())
	}
	for i, v := range choiceList {

		cursor := " " // No cursor used yet
		if a.GopherModel.cursor == i {
			cursor = ">"
		}

		checked := " " // not selected
		if _, ok := a.GopherModel.selected[i]; ok {
			checked = "y"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, v)
	}
	s += "\nPress q to quit.\n"
	return s
}
