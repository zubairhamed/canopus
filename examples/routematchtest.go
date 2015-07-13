package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
	"regexp"
)

func main() {
	e, _ := MatchRoute("rd", "rd")
	FailIfError(e)

	e, _ = MatchRoute("bs", "bs")
	FailIfError(e)

	e, _ = MatchRoute("0/1/2", `^(?P<first>\d+)/(?P<second>\d+)/(?P<third>\d+)$`)
	FailIfError(e)

	e, _ = MatchRoute("0/1/2?abc=123", `^(?P<first>\d+)/(?P<second>\d+)/(?P<third>\d+)\?abc=(?P<fourth>\d+)$`)
	FailIfError(e)

	e, _ = MatchRoute("basic", `^basic$`)
	FailIfError(e)

	e, _ = MatchRoute("0/1/2", `^(?P<obj>\w+)/(?P<inst>\w+)/(?P<rsrc>\w+)$`)
	FailIfError(e)

	re, _ := regexp.Compile(`{[a-z]+}`)
	log.Println(re.FindAllStringSubmatch("{obj}/{inst}/{rsrc}", -1))
}

func FailIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
