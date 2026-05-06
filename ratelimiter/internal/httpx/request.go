package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeJSONBody(
	w http.ResponseWriter,
	r *http.Request,
	dst interface{},
) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	// Prevent unknown fields in request body.
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	// Ensure only a single JSON object exists.
	if decoder.More() {
		return errors.New("multiple JSON objects in request body")
	}

	return nil
}
