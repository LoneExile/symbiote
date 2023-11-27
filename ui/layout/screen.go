package layout

import (
	"fmt"
	"os"
	"os/exec"

	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

type cmdFinishedMsg struct{ err error }

// menuItem represents an item in a menu. It can be a command or lead to a submenu.
type menuItem struct {
	Name    string
	Command func() *exec.Cmd
	SubMenu []menuItem
}

// model represents the application state.
type model struct {
	menus    []menuItem
	cursor   int
	menuPath []int // Stack to keep track of menu navigation
	err      error
}

// Initial model setup with menu items.
func initialModel() model {
	return model{
		menus: []menuItem{
			{
				Name: "AWS",
				SubMenu: []menuItem{
					{Name: "Connect", Command: f.ConnectCmd},
					// {Name: "Connect", Command: a.ConnectToInstance},
					// {Name: "SFTP", Command: a.SFTP},
					// {Name: "List Instances", Command: a.ListInstances},
				},
			},
			// {
			// 	Name: "Local",
			// 	SubMenu: []menuItem{
			// 		{Name: "SSH", Command: test},
			// 		{Name: "SFTP", Command: test},
			// 		{
			// 			Name: "TEST",
			// 			SubMenu: []menuItem{
			// 				{Name: "TEST", Command: test},
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	Name:    "Help",
			// 	Command: test,
			// },
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
				// Navigate into submenu
				m.menuPath = append(m.menuPath, m.cursor)
				m.cursor = 0
			} else if selectedItem.Command != nil {
				// Execute command
				// return m, runCmd(selectedItem.Command)
				if selectedItem.Name == "Help" {
					// Special handling for the Help command
					return m, openHelp() // Implement openHelp() to show help
				} else {
					return m, runCmd(selectedItem.Command)
				}
			}
		case "backspace":
			if len(m.menuPath) > 0 {
				// Navigate up in the menu
				m.menuPath = m.menuPath[:len(m.menuPath)-1]
				m.cursor = 0
			}
		}

		// Handle other message types if necessary

	}
	return m, nil
}

// View renders the UI, which is shown in the terminal.
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

// runCmd is a placeholder for executing commands.
func runCmd(command func() *exec.Cmd) tea.Cmd {
	c := command()
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return cmdFinishedMsg{err}
	})
	// return nil
}

// func runCmd2(command func()) tea.Cmd {
// 	command()
// 	return nil
// }

func test() {

}

// openHelp is a placeholder for opening a help screen.
func openHelp() tea.Cmd {
	// fmt.Println("Showing help")
	// In real application, implement logic to show help screen
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
