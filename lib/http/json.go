package http

import (
	"bytes"
	"encoding/json"
	"io"
)

func MarshalJsonToBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	val, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(val), nil
}

func MarshalJsonToString(body any) (string, error) {
	if body == nil {
		return "", nil
	}

	val, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(val), nil
}
