package goap

import (
	"net"
    "regexp"
    "strings"
)


/*
func GenerateMessageId() uint16 {

}

func GenerateToken() []byte {

}
*/

func SendPacket (msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	b := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	return err
}

func CoreResourcesFromString(str string) []*CoreResource {
    var re = regexp.MustCompile(`(<[^>]+>\s*(;\s*\w+\s*(=\s*(\w+|"([^"\\]*(\\.[^"\\]*)*)")\s*)?)*)`)

    var resources []* CoreResource
    m := re.FindAllString(str, -1)
    for _, match := range m {
        var elemRe = regexp.MustCompile(`<\/[a-zA-Z0-9_%-]+>`)
        elemMatch := elemRe.FindString(match)
        target :=elemMatch[1:len(elemMatch)-1]

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
