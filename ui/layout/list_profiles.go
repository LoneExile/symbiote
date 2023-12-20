package layout

import (
	f "symbiote/cmd/aws/fn"
	// tea "github.com/charmbracelet/bubbletea"
)

func ListProfile() []string {
	l := f.ListProfiles()

	for i, val := range l {
		// if val == "default" {
		// 	l = append(l[:i], l[i+1:]...)
		// }
		if val == "DEFAULT" {
			l = append(l[:i], l[i+1:]...)
		}
	}

	return l
}
