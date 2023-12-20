package fn

import (
	"fmt"
	svc "symbiote/aws"
)

func ListProfiles() []string {
	l, errr := svc.GetAWSProfiles()
	if errr != nil {
		fmt.Println(errr)
	}

	// for _, val := range l {
	// 	fmt.Println(val)
	// }
	return l
}
