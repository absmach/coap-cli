package coap

import (
	"net/url"

	gocoap "github.com/dustin/go-coap"
)

// Option represents CoAP option.
type Option struct {
	ID    gocoap.OptionID
	Value interface{}
}

// ParseOptions parases URL to CoAP options.
func ParseOptions(u *url.URL) []Option {
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
}
