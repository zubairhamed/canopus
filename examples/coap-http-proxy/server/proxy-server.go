package main
import . "github.com/zubairhamed/canopus"

func main() {
	server := NewLocalServer()
	server.SetProxy(PROXY_COAP_HTTP, true)

	server.Start()
}
/*
	- if request contains 'Proxy-Uri' or 'Proxy-Scheme' option with http or https uri
	- if proxy is unwilling, return 5.05 - Proxying Not Supported
	- if timedout, return 5.04 - Gateway Timeout
	- if response is not understood, return 5.02 - Bad Gateway

	- GET
		- Upon success, 2.05 (Content) is returned
		- payload of response MUST be representation of target HTTP resource and Content-Format Option MUST be set
		- include if ETag

		- Client can influence:
			- Accept Option
			- ETag


	SetProxy(HTTP_COAP, bool)
	SetProxy(COAP_COAP, bool)
	SetProxy(COAP_HTTP, bool)

	OnGetRequest
		if allowProxy
			HandleProxy (null or impl)
				if request contains Proxy-Uri or Proxy-Scheme
					if request starts with http or https
						call http request from coap msg
						if  response == bad
							return 5.02 bad gateway
						else
						if response == timeout
							return 5.04 gateway timeout
						else
							if response type != accept type
								return 5.02 bad gateway

							if method == GET
								return 2.05 (Content)
							if method == PUT
								return 2.04 (Changed)
							if method == DELETE
								return 2.02 Deleted
							if method == POST
								return 2.01 (Created)
							return coap msg
					else
					if request starts with coap or coaps
		else
			throw 5.05 not supported


 */
