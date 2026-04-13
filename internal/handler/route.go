package handler

import "net/http"

type Route struct {
	Handler http.HandlerFunc
	Method  string
	Path    string
}
