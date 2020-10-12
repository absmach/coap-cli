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
		{
			ID:    gocoap.URIHost,
			Value: u.Hostname(),
		},
		{
			ID:    gocoap.URIPath,
			Value: u.Path[1:],
		},
		{
			ID:    gocoap.URIPort,
			Value: u.Port(),
		},
		{
			ID:    gocoap.URIQuery,
			Value: u.RawQuery,
		},
	}
}
