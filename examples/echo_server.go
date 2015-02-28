package main

import (
	"fmt"
)

type Abc struct {

}

func main() {
	mymap := make(map[string] Abc)

	fmt.Println(mymap["abc"])
}

