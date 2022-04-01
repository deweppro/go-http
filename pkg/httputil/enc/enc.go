package enc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		Error(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b) //nolint: errcheck
}

func XML(w http.ResponseWriter, v interface{}) {
	b, err := xml.Marshal(v)
	if err != nil {
		Error(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b) //nolint: errcheck
}

func Error(w http.ResponseWriter, v error) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(v.Error())) //nolint: errcheck
}

func Stream(w http.ResponseWriter, v []byte, filename string) {
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	w.Write(v) //nolint: errcheck
}

func Raw(w http.ResponseWriter, v []byte) {
	w.Header().Add("Content-Type", http.DetectContentType(v))
	w.WriteHeader(http.StatusOK)
	w.Write(v) //nolint: errcheck
}
