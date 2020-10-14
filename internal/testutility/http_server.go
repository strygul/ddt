package testutility

import (
	"net/http"
	"sync"
)

type MyHandler struct {
	sync.Mutex
	count int
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// create response binary data
	data := []byte(`{"foo": "bar"}`) // slice of bytes
	// write `data` to response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
