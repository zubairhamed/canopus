package main

import . "github.com/zubairhamed/goap"

func main() {
    MatchRoute("rd", "rd")
    MatchRoute("bs", "bs")

    MatchRoute("x")
}


/*


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