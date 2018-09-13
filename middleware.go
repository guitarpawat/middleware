package middleware

import (
	"net/http"
)

// ValueMap sends context as a key-value between the chaining of Middleware.
type ValueMap map[string]interface{}

// Doable interface for wrapping http handler.
type Doable interface {
	Do(w http.ResponseWriter, r *http.Request, v *ValueMap)
}

// Middleware struct is implements http.Handler.
type Middleware struct {
	This Doable
	Next *Middleware
	*ValueMap
}

// ServeHTTP is implementation of http.Handler. It creates a new ValueMap if nil,
// runs the wrapped Doable, and calls the next Middleware if ShouldNext is true
// and next Middleware is not nil.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.ValueMap == nil {
		m.ValueMap = &(ValueMap{})
	}
	m.This.Do(w, r, m.ValueMap)
	shouldNext, ok := (*m.ValueMap)["next"].(bool)
	if ok && shouldNext && m.Next != nil {
		m.Next.ValueMap = m.ValueMap
		m.Next.ServeHTTP(w, r)
	}
}

// MakeMiddleware creates a new Middleware with chaining Doable stuffs.
func MakeMiddleware(initial *ValueMap, stuff ...Doable) Middleware {
	switch len(stuff) {
	case 0:
		return Middleware{}
	case 1:
		return Middleware{
			This:     stuff[0],
			Next:     nil,
			ValueMap: initial,
		}
	default:
		nextMiddleware := MakeMiddleware(nil, stuff[1:]...)
		return Middleware{
			This:     stuff[0],
			Next:     &nextMiddleware,
			ValueMap: initial,
		}
	}
}
