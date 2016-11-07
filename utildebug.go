package canopus

import (
	"log"
)

// PrintOptions pretty prints out a given Message's options
func PrintOptions(msg Message) {
	opts := msg.GetAllOptions()
	log.Println(" - - - OPTIONS - - - ")
	if len(opts) > 0 {
		for _, opts := range msg.GetAllOptions() {
			log.Println("Code/Number: ", opts.GetCode(), ", Name: ", OptionNumberToString(opts.GetCode()), ", Value: ", opts.GetValue())
		}
	} else {
		log.Println("None")
	}
}

// PrintMessage pretty prints out a given Message
func PrintMessage(msg Message) {
	log.Println("= = = = = = = = = = = = = = = = ")
	log.Println("Code: ", msg.GetCode())
	log.Println("Code String: ", CoapCodeToString(msg.GetCode()))
	log.Println("MessageId: ", msg.GetMessageId())
	log.Println("MessageType: ", msg.GetMessageType())
	log.Println("Token: ", string(msg.GetToken()))
	log.Println("Token Length: ", msg.GetTokenLength())
	log.Println("Payload: ", payload.PayloadAsString(msg.GetPayload()))
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
