package layout

import (
	"os"
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

func ForwardDB(m model) tea.Cmd {
	port := m.TextForm.Inputs[0].Value() + ":" + m.TextForm.Inputs[1].Value()
	c := f.RdsCmd(port, m.CurrentP)
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
