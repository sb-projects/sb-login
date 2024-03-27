package response

import (
	"encoding/json"
	"net/http"
)

type (
	ErrResp struct {
		Error string `json:"error,omitempty"`
		Data  any    `json:"data,omitempty"`
	}
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONError(w http.ResponseWriter, status int, data error) error {
	return JSONWithHeaders(w, status, ErrResp{Error: data.Error()}, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
