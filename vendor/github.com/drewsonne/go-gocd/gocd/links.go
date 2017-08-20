package gocd

import (
	"fmt"
	"net/url"
)

type linkField map[string]map[string]string
type linkHref struct {
	H string `json:"href"`
}

func unmarshallLinkField(d linkField, field string, destination **url.URL) error {
	var e error
	if h, ok := d[field]["href"]; ok && h != "" {
		*destination, e = url.Parse(h)
		return e
	}
	return fmt.Errorf("'%s' was not present in `%s`", field, d)
}
