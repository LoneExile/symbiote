package layout

import (
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

func ConnectCmd(m model) tea.Cmd {
	c := f.ConnectCmd()

	teaCmds := tea.ExecProcess(c, func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})

	if true {
		runEchoCmd := echoCmd(c)
		return runCmds([]tea.Cmd{runEchoCmd, teaCmds})
	}

	return teaCmds
}
