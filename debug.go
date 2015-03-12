package goap

import "log"


func PrintOptions(msg *Message) {
	log.Println("==== PRINT OPTIONS ===")

	for _, opts := range msg.Options {
		log.Println("Code/Number: ", opts.Code)
		log.Println("Name: ", OptionNumberToString(opts.Code))
		log.Println("Value: ", opts.Value)

		log.Println("---")
	}

	log.Println("======================")
}

func OptionNumberToString(o OptionCode) string {
	switch o {
	case OPTION_IF_MATCH:
		return "If-Match"

	case OPTION_URI_HOST:
		return "Uri-Host"

	case OPTION_ETAG:
		return "ETag"

	case OPTION_IF_NONE_MATCH:
		return "If-None-Match"

	case OPTION_URI_PORT:
		return "Uri-Port"

	case OPTION_LOCATION_PATH:
		return "Location-Path"

	case OPTION_URI_PATH:
		return "Uri-Path"

	case OPTION_CONTENT_FORMAT:
		return "Content-Format"

	case OPTION_MAX_AGE:
		return "Max-Age"

	case OPTION_URI_QUERY:
		return "Uri-Query"

	case OPTION_ACCEPT:
		return "Accept"

	case OPTION_LOCATION_QUERY:
		return "Location-Query"

	case OPTION_BLOCK2:
		return "Block2"

	case OPTION_BLOCK1:
		return "Block1"

	case OPTION_PROXY_URI:
		return "Proxy-Uri"

	case OPTION_PROXY_SCHEME:
		return "Proxy-Scheme"

	case OPTION_SIZE1:
		return "Size1"

	default:
		return ""
	}
	return ""
}


/*

 */
