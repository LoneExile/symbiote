package layout

import (
	f "symbiote/cmd/aws/fn"
)

func ListProfile() []string {
	l := f.ListProfiles()

	for i, val := range l {
		if val == "DEFAULT" {
			l = append(l[:i], l[i+1:]...)
		}
	}

	return l
}
