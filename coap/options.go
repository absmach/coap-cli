package coap

import (
	"fmt"
	"net/url"

	gocoap "github.com/dustin/go-coap"
)

// Option represents CoAP option.
type Option struct {
	ID    gocoap.OptionID
	Value interface{}
}

func ParseOptions(u *url.URL) []Option {
	// var ret []Option
	fmt.Println(u.Path)
	return []Option{
		Option{
			ID:    gocoap.URIHost,
			Value: u.Hostname(),
		},
		Option{
			ID:    gocoap.URIPath,
			Value: u.Path,
		},
		Option{
			ID:    gocoap.URIPort,
			Value: u.Port(),
		},
		Option{
			ID:    gocoap.URIQuery,
			Value: u.RawQuery,
		},
	}
	// vals := u.Query()
	// for k, v := range vals {
	// 	fmt.Println("K:", k, "V:", v)
	// 	// o := Option{ID: gocoap.URIPath}
	// }
	// return ret
}
