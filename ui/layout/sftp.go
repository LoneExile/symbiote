package layout

import (
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
)

func EicSFTPCmd(m model) tea.Cmd {
	tunnelCmd := f.EicSFTPCmd()
	eicCmd := runBgCmd(tunnelCmd, m.CurrentCmd.Wording)

	if true {
		runEchoCmd := echoCmd(tunnelCmd)
		return runCmds([]tea.Cmd{runEchoCmd, eicCmd})
	}

	return eicCmd
}

func SFTPConnectCmd(m model) tea.Cmd {
	return tea.ExecProcess(f.SFTPConnectCmd(), func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})
}
