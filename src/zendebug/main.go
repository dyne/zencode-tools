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
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	initialInputs = 3
	maxInputs     = 6
	minInputs     = 2
	helpHeight    = 5
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type keymap = struct {
	next, prev, /* add, remove,  */ exec, quit key.Binding
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = " "
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type model struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int
}

func newModel() model {
	m := model{
		inputs: make([]textarea.Model, initialInputs),
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("ctrl+right"),
				key.WithHelp("ctrl ->", "focus right"),
			),
			prev: key.NewBinding(
				key.WithKeys("ctrl+left"),
				key.WithHelp("ctrl <-", "focus left"),
			),
			exec: key.NewBinding(
				key.WithKeys("ctrl+down"),
				key.WithHelp("ctrl down", "EXEC"),
			),
			// add: key.NewBinding(
			// 	key.WithKeys("ctrl+n"),
			// 	key.WithHelp("ctrl+n", "add an editor"),
			// ),
			// remove: key.NewBinding(
			// 	key.WithKeys("ctrl+w"),
			// 	key.WithHelp("ctrl+w", "remove an editor"),
			// ),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	for i := 0; i < initialInputs; i++ {
		m.inputs[i] = newTextarea()
	}
	m.inputs[m.focus].Focus()
	m.updateKeybindings()
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)

			// TODO: not working
		case key.Matches(msg, m.keymap.exec):
			//m.inputs[1].SetValue( m.inputs[0].Value() )
			//cmd := tea.ExecProcess( exec.Command("zenroom", m.inputs[0].Value()) ,
			//	func(err error) tea.Msg {
			//		m.inputs[2].SetValue(err.Error())
			//		return err.Error()
			//	})
			//cmds = append(cmds, cmd)
			dataFile, _ := os.Create("tempData")
			defer dataFile.Close()
			dataFile.WriteString(m.inputs[1].Value())
			dataFile.Sync()

			scriptFile, _ := os.Create("tempScript")
			defer scriptFile.Close()
			scriptFile.WriteString(m.inputs[0].Value())
			scriptFile.Sync()
			cmd := exec.Command("zenroom","-z",scriptFile.Name(),"-a",dataFile.Name())
			out, err := cmd.CombinedOutput()
			if err != nil {
				m.inputs[2].SetValue( string(out) + err.Error() )
			} else {
				m.inputs[2].SetValue( string(out) )
			}
			os.Remove(dataFile.Name())
			os.Remove(scriptFile.Name())
		// case key.Matches(msg, m.keymap.add):
		// 	m.inputs = append(m.inputs, newTextarea())
		// case key.Matches(msg, m.keymap.remove):
		// 	m.inputs = m.inputs[:len(m.inputs)-1]
		// 	if m.focus > len(m.inputs)-1 {
		// 		m.focus = len(m.inputs) - 1
		// 	}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.updateKeybindings()
	m.sizeInputs()

	// Update all textareas
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

func (m *model) updateKeybindings() {
	// m.keymap.exec.SetEnabled(false)
	// m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	// m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		// m.keymap.add,
		// m.keymap.remove,
		m.keymap.exec,
		m.keymap.quit,
	})

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
