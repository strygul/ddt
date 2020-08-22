package entity

import "net/http"

type Executor struct {
	Path Path
	Client http.Client
}
