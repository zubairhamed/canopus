package main

import (
	. "github.com/zubairhamed/goap"
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

/*
    OnNewRoute
        Get all values between #{ }
        Construct New RegEx
            Create SubGroups
            Escape any RegEx Values
        Compile and Store Compiled RegEx



(?P<first>\d+)\.(\d+).(?P<second>\d+)`

/bs 											POST
/rd												POST

/{obj}/{inst}/{rsrc}			PUT				// Write & Write Attribute

/{obj}/{inst}/{rsrc}			GET
/{obj}/{inst}/{rsrc}			GET 			application/link-format
/{obj}/{inst}/{rsrc}			GET + Observe Options		// Observe

/{obj}/{inst}/{rsrc}			PUT/POST	// Write and Execute


/{obj}/{inst}/{rsrc}			Reset		// Cancel Observation
/{obj}/{inst}/{rsrc}			Async Response	// Notify

/{obj}/{inst}							DELETE		// Delete
/{obj}/{inst}							POST			// Create

/{loc}										PUT
/{loc}										DELETE


*/
