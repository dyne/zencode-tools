package main
import (
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)
type ZencodeListModel struct {
	list list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}
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

type ZencodeStatement struct {
	scenario  string
	statement string
}

func (i ZencodeStatement) Title()   string { return i.statement }
func (i ZencodeStatement) Description()    string { return i.scenario }
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
		statements = append(statements, ZencodeStatement {
			scenario: "",
			statement: "Given I " + z.Given[i],
		})
	}
	for k, v := range z.When {
		for i := 0; i < len(v); i++ {
			var scenario = ""
			if k != "default" {
				scenario = k
			}
			statements = append(statements, ZencodeStatement {
				scenario: scenario,
				statement: "When I " + v[i],
			})
		}
	}
	for i := 0; i < len(z.Then); i++ {
		statements = append(statements, ZencodeStatement {
			scenario: "",
			statement: "Then I " + z.Then[i],
		})
	}

	sort.Sort(ByFilterValue(statements))


	return statements
}

func newZencodeListModel(zen ZenStatements) ZencodeListModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Make initial list of items
	items := createKeyValueList(zen)

	// Setup list
	delegate := newItemDelegate(delegateKeys)
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

	return ZencodeListModel{
		list:          zencodeList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
	}
}


func (m ZencodeListModel) Update(msg tea.Msg) (ZencodeListModel, tea.Cmd) {
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

func (m ZencodeListModel) View() string {
	return m.list.View()
}
