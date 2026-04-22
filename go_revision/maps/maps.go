package maps

import (
	"fmt"
)

func MapsExample() {
	ages := map[string]int{
		"Alice": 30,
		"Bob": 10,
	}

	fmt.Printf("ages map: %#v\n", ages)

	is_druggie := make(map[string]bool, 10)	// here 10 is capacity and not length, we can't have zero values for maps hence we can only tell the runtime about the capacity for maps which would prevent re-copies and hence re-allocation of memory
	is_druggie["Alice"] = true
	fmt.Printf("is_druggie map: %#v\n", is_druggie)
}
