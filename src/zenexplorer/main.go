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

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"net"
	"path/filepath"
	"sort"
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

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	zencodeItems *ZenStatements
	keys         *listKeyMap
	delegateKeys *delegateKeyMap

	serverStarted bool
}

type ZencodeStatement struct {
	scenario  string
	statement string
}

func (i ZencodeStatement) Title() string       { return i.statement }
func (i ZencodeStatement) Description() string { return i.scenario }
func (i ZencodeStatement) FilterValue() string { return i.statement }

type ByFilterValue []list.Item

func (a ByFilterValue) Len() int           { return len(a) }
func (a ByFilterValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFilterValue) Less(i, j int) bool { return a[i].FilterValue() < a[j].FilterValue() }

func createKeyValueList(z ZenStatements) []list.Item {
	if z.mtx == nil {
		z.reset()
	}

	var statements []list.Item

	z.mtx.Lock()
	defer z.mtx.Unlock()

	for i := 0; i < len(z.Given); i++ {
		statements = append(statements, ZencodeStatement{
			scenario:  "",
			statement: "Given I " + z.Given[i],
		})
	}
	for k, v := range z.When {
		for i := 0; i < len(v); i++ {
			var scenario = ""
			if k != "default" {
				scenario = k
			}
			statements = append(statements, ZencodeStatement{
				scenario:  scenario,
				statement: "When I " + v[i],
			})
		}
	}
	for i := 0; i < len(z.Then); i++ {
		statements = append(statements, ZencodeStatement{
			scenario:  "",
			statement: "Then I " + z.Then[i],
		})
	}

	sort.Sort(ByFilterValue(statements))

	return statements
}
func newModel(sock net.Listener) model {
	var (
		zencodeItems   ZenStatements
		delegateKeys               = newDelegateKeyMap()
		listKeys                   = newListKeyMap()
		chanStatements chan string = nil
	)

	serverStarted := false
	if sock != nil {
		serverStarted = true
		chanStatements = make(chan string)
		go startServer(sock, chanStatements)
	}

	// Make initial list of items
	items := createKeyValueList(zencodeItems)

	// Setup list
	delegate := newItemDelegate(delegateKeys, chanStatements)
	zencodeList := list.New(items, delegate, 0, 0)
	zencodeList.Title = "Statements"
	zencodeList.Styles.Title = titleStyle
	zencodeList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:          zencodeList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		zencodeItems:  &zencodeItems,
		serverStarted: serverStarted,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {
	p, err := filepath.Abs(filepath.Join("/", "tmp", "zencode-tools"))
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}
	socketFile := filepath.Join(p, "explorer.sock")
	l, err := net.Listen("unix", socketFile)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	if err := tea.NewProgram(newModel(l)).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	os.Remove(socketFile)
}
