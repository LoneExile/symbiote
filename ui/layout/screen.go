package layout

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	Name    string
	Command func() tea.Cmd
	SubCmd  func() tea.Cmd
	SubMenu []menuItem
}

type model struct {
	menus    []menuItem
	cursor   int
	menuPath []int
	err      error
}

// type commandStartedMsg struct{}
type FoundSubCmd struct{}
type commandCompletedMsg struct{ err error }
type commandFailedMsg struct{ err error }

// Initial model setup with menu items.
func initialModel() model {
	return model{
		menus: []menuItem{
			{
				Name: "AWS",
				SubMenu: []menuItem{
					{Name: "Connect", Command: ConnectCmd, SubCmd: nil},
					{Name: "SFTP", Command: EicSFTPCmd, SubCmd: SFTPConnectCmd},
				},
			},
		},
		menuPath: make([]int, 0),
	}
}

// Init is called when the program starts.
func (m model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.getCurrentMenu())-1 {
				m.cursor++
			}
		case "enter", " ":
			currentMenu := m.getCurrentMenu()
			selectedItem := currentMenu[m.cursor]

			if len(selectedItem.SubMenu) > 0 {
				m.menuPath = append(m.menuPath, m.cursor)
				m.cursor = 0
			} else if selectedItem.Command != nil {
				if selectedItem.Name == "Help" {
					return m, openHelp()
				} else {
					return m, selectedItem.Command()
				}
			}
		case "backspace":
			if len(m.menuPath) > 0 {
				// Navigate up in the menu
				m.menuPath = m.menuPath[:len(m.menuPath)-1]
				m.cursor = 0
			}
		}
	}

	switch msg := msg.(type) {
	case FoundSubCmd:
		currentMenu := m.getCurrentMenu()
		selectedItem := currentMenu[m.cursor]
		if selectedItem.SubCmd != nil {
			return m, selectedItem.SubCmd()
		}
	case commandCompletedMsg:
		if msg.err != nil {
			fmt.Println("Command completed with error:", msg.err)
		} else {
			fmt.Println("Command completed successfully")
		}
	case commandFailedMsg:
		if msg.err != nil {
			fmt.Println("Command failed with error:", msg.err)
		} else {
			fmt.Println("Command failed")
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "\n\n"
	currentMenu := m.getCurrentMenu()

	for i, item := range currentMenu {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s  %s\n", cursor, item.Name)
	}
	s += "\nPress q to quit, backspace to go back.\n"
	return s
}

// getCurrentMenu returns the current menu based on the navigation path.
func (m *model) getCurrentMenu() []menuItem {
	menu := m.menus
	for _, idx := range m.menuPath {
		menu = menu[idx].SubMenu
	}
	return menu
}

func openHelp() tea.Cmd {
	// fmt.Println("Showing help")
	return nil
}

// main function to run the Bubble Tea program.
func Screen() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
