package main
import "github.com/zubairhamed/canopus"

func main() {

	type ProxyPass struct {
		in string
		out string
	}

	var server *canopus.CoapServer
	server.AddProxyPass(in, out)
	server.RemoveProxyPass(in)

	server.AddProxyPassReverse(in, out)
	server.AddRemoveProxyPassReverse


}
