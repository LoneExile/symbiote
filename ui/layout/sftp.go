package layout

import (
	"bufio"
	"strings"
	f "symbiote/cmd/aws/fn"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/creack/pty"
)

func EicSFTPCmd() tea.Cmd {
	return func() tea.Msg {
		tunnelCmd := f.EicSFTPCmd()
		ptmx, err := pty.Start(tunnelCmd)
		if err != nil {
			return commandFailedMsg{err}
		}

		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Listening") {
				return FoundSubCmd{}
			}
		}
		return commandFailedMsg{err}
	}
}

func SFTPConnectCmd() tea.Cmd {
	return tea.ExecProcess(f.SFTPConnectCmd(), func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})
}
