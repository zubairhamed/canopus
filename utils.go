package goap

import (
	"net"
	"regexp"
	"strings"
)

var MESSAGEID_CURR = 0

/*
func GenerateMessageId() uint16 {

}

func GenerateToken() []byte {

}
*/

// Sends a CoAP Message to UDP address
func SendPacket(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	b := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	return err
}

// Converts to CoRE Resources Object from a CoRE String
func CoreResourcesFromString(str string) []*CoreResource {
	var re = regexp.MustCompile(`(<[^>]+>\s*(;\s*\w+\s*(=\s*(\w+|"([^"\\]*(\\.[^"\\]*)*)")\s*)?)*)`)
	var elemRe = regexp.MustCompile(`<\/[a-zA-Z0-9_%-]+>`)

	var resources []*CoreResource
	m := re.FindAllString(str, -1)

	for _, match := range m {
		elemMatch := elemRe.FindString(match)
		target := elemMatch[1 : len(elemMatch)-1]

		resource := NewCoreResource()
		resource.Target = target

		attrs := strings.Split(match[len(elemMatch)+1:], ";")

		for _, attr := range attrs {
			pair := strings.Split(attr, "=")

			resource.AddAttribute(pair[0], strings.Replace(pair[1], "\"", "", -1))
		}

		resources = append(resources, resource)
	}
	return resources
}
