package layout

import (
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

func ConnectCmd() tea.Cmd {
	return tea.ExecProcess(f.ConnectCmd(), func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})
}
