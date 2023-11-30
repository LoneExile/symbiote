package layout

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	Name    string
	Command []Cmd
	SubMenu []menuItem
}

type Cmd struct {
	Cmd     func(model) tea.Cmd
	Type    string
	Wording string
}

type currentCmd struct {
	stage   int
	Wording string
	Type    string
}

type model struct {
	Menus      []menuItem
	Cursor     int
	MenuPath   []int
	CurrentCmd currentCmd
}

// type commandStartedMsg struct{}
type FoundSubCmd struct{}
type commandCompletedMsg struct{ err error }
type commandFailedMsg struct{ err error }

// Initial model setup with menu items.
func initialModel() model {
	return model{
		Menus: []menuItem{
			{
				Name: "AWS",
				SubMenu: []menuItem{
					{Name: "Connect", Command: []Cmd{
						{Cmd: ConnectCmd, Type: "exec"},
					},
					},
					{
						Name: "SFTP",
						Command: []Cmd{
							{Cmd: EicSFTPCmd, Type: "bg", Wording: "Listening"},
							{Cmd: SFTPConnectCmd, Type: "exec"},
						},
					}, // ,
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
		MenuPath: make([]int, 0),
	}
}

// Init is called when the program starts.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.getCurrentMenu())-1 {
				m.Cursor++
			}
		case "enter", " ":
			currentMenu := m.getCurrentMenu()
			selectedItem := currentMenu[m.Cursor]

			if len(selectedItem.SubMenu) > 0 {
				m.MenuPath = append(m.MenuPath, m.Cursor)
				m.Cursor = 0
			} else if selectedItem.Command != nil {
				if selectedItem.Name == "Help" {
					return m, openHelp()
				} else {
					m.CurrentCmd.stage = 0
					m.CurrentCmd.Wording = selectedItem.Command[0].Wording
					m.CurrentCmd.Type = selectedItem.Command[0].Type
					return m, selectedItem.Command[0].Cmd(m)
				}
			}
		case "backspace":
			if len(m.MenuPath) > 0 {
				// Navigate up in the menu
				m.MenuPath = m.MenuPath[:len(m.MenuPath)-1]
				// TODO: set cursor to the last selected item
				m.Cursor = 0
			}
		}
	}

	switch msg := msg.(type) {
	case FoundSubCmd:
		currentMenu := m.getCurrentMenu()
		selectedItem := currentMenu[m.Cursor]
		// if selectedItem.Command[1].Cmd != nil {
		// 	return m, selectedItem.Command[1].Cmd()
		// }

		stage := m.CurrentCmd.stage
		if stage < len(selectedItem.Command)-1 {
			m.CurrentCmd.stage++
			m.CurrentCmd.Wording = selectedItem.Command[stage+1].Wording
			m.CurrentCmd.Type = selectedItem.Command[stage+1].Type
			return m, selectedItem.Command[stage+1].Cmd(m)
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
		if m.Cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s  %s\n", cursor, item.Name)
	}
	s += "\nPress q to quit, backspace to go back.\n"
	return s
}

// getCurrentMenu returns the current menu based on the navigation path.
func (m *model) getCurrentMenu() []menuItem {
	menu := m.Menus
	for _, idx := range m.MenuPath {
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
