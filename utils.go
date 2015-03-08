package goap

import (
	"regexp"
	"strings"
	"math/rand"
	"time"
)

var MESSAGEID_CURR = 0

/*
func GenerateMessageId() uint16 {

}
*/

var genChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
func GenerateToken (l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	token := make([]rune, l)
	for i := range token {
		token[i] = genChars[rand.Intn(len(genChars))]
	}
	return string(token)
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
