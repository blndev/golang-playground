package main

import (
	"fmt"
)

type tags map[string]string

func main() {
	fmt.Println("Map - Playground")

	var m = tags{} //{"test", "value"}
	fmt.Printf("%#v\n", m)

	m["hello"] = "world"
	fmt.Printf("%#v\n", m)

	if val, ok := m["hello"]; ok {
		//do something here
		fmt.Println("found %s in map for key %s", val, "hello")
	}

}
