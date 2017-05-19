package cmap_test

import (
	"fmt"

	"github.com/corsc/go-commons/concurrency/cmap"
)

func ExampleMap_defaultSharding() {
	myMap := cmap.New()

	val, err := myMap.Get("key")
	fmt.Printf("Value/Err: %v/%v\n", val, err)

	err = myMap.Set("key", "foo")
	fmt.Printf("Err: %v\n", err)

	val, err = myMap.Get("key")
	fmt.Printf("Value/Err: %v/%v\n", val, err)

	// Output:
	// Value/Err: <nil>/no such item
	// Err: <nil>
	// Value/Err: foo/<nil>
}
