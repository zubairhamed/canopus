package canopus

import (
	"log"
)

// PrintOptions pretty prints out a given Message's options
func PrintOptions(msg *Message) {
	opts := msg.Options
	log.Println(" - - - OPTIONS - - - ")
	if len(opts) > 0 {
		for _, opts := range msg.Options {
			log.Println("Code/Number: ", opts.GetCode(), ", Name: ", OptionNumberToString(opts.GetCode()), ", Value: ", opts.GetValue())
		}
	} else {
		log.Println("None")
	}
}

// PrintMessage pretty prints out a given Message
func PrintMessage(msg *Message) {
	log.Println("= = = = = = = = = = = = = = = = ")
	log.Println("Code: ", msg.Code)
	log.Println("Code String: ", CoapCodeToString(msg.Code))
	log.Println("MessageId: ", msg.MessageID)
	log.Println("MessageType: ", msg.MessageType)
	log.Println("Token: ", string(msg.Token))
	log.Println("Token Length: ", msg.GetTokenLength())
	log.Println("Payload: ", PayloadAsString(msg.Payload))
	PrintOptions(msg)
	log.Println("= = = = = = = = = = = = = = = = ")

}

// OptionNumberToString returns the string representation of a given Option Code
func OptionNumberToString(o OptionCode) string {
	switch o {
	case OptionIfMatch:
		return "If-Match"

	case OptionURIHost:
		return "Uri-Host"

	case OptionEtag:
		return "ETag"

	case OptionIfNoneMatch:
		return "If-None-Match"

	case OptionURIPort:
		return "Uri-Port"

	case OptionLocationPath:
		return "Location-Path"

	case OptionURIPath:
		return "Uri-Path"

	case OptionContentFormat:
		return "Content-Format"

	case OptionMaxAge:
		return "Max-Age"

	case OptionURIQuery:
		return "Uri-Query"

	case OptionAccept:
		return "Accept"

	case OptionLocationQuery:
		return "Location-Query"

	case OptionBlock2:
		return "Block2"

	case OptionBlock1:
		return "Block1"

	case OptionProxyURI:
		return "Proxy-Uri"

	case OptionProxyScheme:
		return "Proxy-Scheme"

	case OptionSize1:
		return "Size1"

	case OptionSize2:
		return "Size2"

	default:
		return ""
	}
}
