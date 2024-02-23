package layout

import (
	"os"
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

func ForwardDB(m model) tea.Cmd {
	dbInstances := f.ListDBInstances(m.CurrentP)

	if len(dbInstances) == 0 {
		// os.Exit(1)
		return func() tea.Msg {
			return commandFailedMsg{err: nil}
		}
	}

	return func() tea.Msg {
		return startDBSelectionMsg{DBInstances: dbInstances}
	}

}

func connectToDB(m model, port string) tea.Cmd {
	c := f.RdsCmdSelected(port, m.CurrentP, m.SelectedDB)
	if c == nil {
		os.Exit(1)
	}

	teaCmds := tea.ExecProcess(c, func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})

	if true {
		runEchoCmd := echoCmd(c)
		return runCmds([]tea.Cmd{runEchoCmd, teaCmds})
	}
	return teaCmds
}
