package main

import (
	. "github.com/zubairhamed/goap"
)

/*
Forward Proxy
USAGE: Server().AllowForwarding(bool)

- if connection from proxy times out, send to Client 5.04 Gateway Timeout
- if proxy gets error, return 5.02 Bad Gateway
- else return response to client
- if allowForwarding = false, return 5.05 (Proxying Not Supported)
- host and port = authority/target, treat it as local (non-proxied)

- Client Set Proxy-Uri options
- Set Uri-host, Uri-port, Uri-path, Uri-Query (proxy splits into individual options)
- Alternatively: can use Proxy-Scheme option (?)
 */

func main() {
	originServer := NewLocalServer()

	go originServer.Start()

	forwardProxy := NewLocalServer()
	forwardProxy.AllowForwarding(true)

	go forwardProxy.Start()

	client := NewClient()
	defer client.Close()

	// Create a forward proxy request and dispatch
}
