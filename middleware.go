package middleware

import (
	"net/http"
)

// Doable interface for wrapping http handler.
type Doable interface {
	Do(w http.ResponseWriter, r *http.Request, v *ValueMap)
}

// DoableFunc is a wrapper for Doable interface.
type DoableFunc func(w http.ResponseWriter, r *http.Request, v *ValueMap)

// Do is a helper method.
func (f DoableFunc) Do(w http.ResponseWriter, r *http.Request, v *ValueMap) {
	f(w, r, v)
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
	shouldNext, ok := m.Get("next").(bool)
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
	default:
		return *NewMiddleware(initial, stuff...)
	}
}

// NewMiddleware creates a new Middleware pointer's with chaining Doable stuffs.
func NewMiddleware(initial *ValueMap, stuff ...Doable) *Middleware {
	switch len(stuff) {
	case 0:
		return nil
	default:
		return &Middleware{
			This:     stuff[0],
			Next:     NewMiddleware(nil, stuff[1:]...),
			ValueMap: initial,
		}
	}
}
