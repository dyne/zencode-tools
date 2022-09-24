/* Software Tools to work with Zenroom (https://dev.zenroom.org)
 *
 * Copyright (C) 2022 Dyne.org foundation
 * Originally written as example code in Bubblewrap
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type Window int

const (
        ZencodeListWindow Window = iota
		SelectScenarioWindow
)


func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.currentWindow = ZencodeListWindow
			return m, nil
		case "2":
			m.currentWindow = SelectScenarioWindow
			return m, nil
		}
	}
	switch(m.currentWindow) {
	case ZencodeListWindow:
		newListModel, cmds := m.zencodeList.Update(msg)
		m.zencodeList = newListModel
		return m, cmds
	case SelectScenarioWindow:
		newListModel, cmds := m.selectScenario.Update(msg)
		m.selectScenario = newListModel
		return m, cmds
	}
	return m, nil
}

func (m model) View() string {
	var display string
	switch(m.currentWindow) {
	case ZencodeListWindow:
		display = m.zencodeList.View()
	case SelectScenarioWindow:
		display = m.selectScenario.View()
	}
	return appStyle.Render(display)
}

type model struct {
	currentWindow Window
	zencodeList  ZencodeListModel
	selectScenario SelectScenarioModel
	zencodeItems *ZenStatements
}

func newModel() model {
	var (
		zen ZenStatements
	)

	zencodeList := newZencodeListModel(zen)
	selectScenario := newSelectScenarioModel(zen)

	return model{
		currentWindow: ZencodeListWindow,
		zencodeList:   zencodeList,
		selectScenario: selectScenario,
		zencodeItems:  &zen,
	}
}


func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
