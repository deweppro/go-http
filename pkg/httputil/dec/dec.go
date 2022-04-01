package dec

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/deweppro/go-http/internal"
)

func JSON(r *http.Request, v interface{}) error {
	b, err := internal.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func XML(r *http.Request, v interface{}) error {
	b, err := internal.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, v)
}
