package procspy_test

import (
	"fmt"

	"github.com/john-pierce/procspy"
)

func Example() {
	lookupProcesses := true
	cs, err := procspy.Connections(lookupProcesses)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TCP Connections:\n")
	for c := cs.Next(); c != nil; c = cs.Next() {
		if !procspy.IsListening(*c) {
			continue
		}
		fmt.Printf(" - %v\n", c)
	}
}
