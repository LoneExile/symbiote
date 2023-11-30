package layout

import (
	"bufio"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/creack/pty"
)

func runCmds(teaCmds []tea.Cmd) tea.Cmd {
	teaCmds = append(teaCmds, func() tea.Msg {
		return commandCompletedMsg{nil}
	})
	return tea.Sequence(teaCmds...)
}

func echoCmd(cmd *exec.Cmd) tea.Cmd {
	cmdStr := "\n" + cmd.String() + "\n"
	echoCmd := exec.Command("echo", cmdStr)

	return tea.ExecProcess(echoCmd, func(err error) tea.Msg {
		return commandCompletedMsg{err}
	})
}

func runBgCmd(cmd *exec.Cmd, word string) tea.Cmd {
	return func() tea.Msg {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return commandFailedMsg{err}
		}
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, word) {
				return FoundSubCmd{}
			}
		}
		return commandFailedMsg{err}
	}
}
