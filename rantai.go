// rantai is an http middleware chainer that provides easy way to chains middlewares.
// It provides some features to manage the chainer such as Extend and Exclude.
package rantai

import (
	"net/http"

	"github.com/lutffmn/kmpr"
)

type Middlewares []func(http.Handler) http.Handler

// Rantai is a struct to keep track of Middlewares
type Rantai struct {
	Middlewares Middlewares
}

// New is going to create a new instance of Rantai struct
func New(middlewares ...func(http.Handler) http.Handler) *Rantai {
	return &Rantai{
		Middlewares: middlewares,
	}
}

// Chain will execute the chaining process, it takes an http.Handler as arg
func (r *Rantai) Chain(handler http.Handler) http.Handler {
	if handler == nil {
		panic("handler couldn't be nil")
	}

	for _, m := range r.Middlewares {
		handler = m(handler)
	}

	return handler
}

// ChainF will execute the chaining process, it takes func(http.ResponseWriter, *http.Request) as arg
func (r *Rantai) ChainF(handlerFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	if handlerFunc == nil {
		panic("handlerFunc couldn't be nil")
	}
	return r.Chain(http.HandlerFunc(handlerFunc))
}

// Extend will return a new Rantai instance with extended Middlewares
func (r *Rantai) Extend(middleware ...func(http.Handler) http.Handler) *Rantai {
	newMiddlewares := r.Middlewares
	for _, m := range middleware {
		newMiddlewares = append(newMiddlewares, m)
	}

	return &Rantai{
		Middlewares: newMiddlewares,
	}
}

// Exclude will return a new Rantai instance with excluded middleware that passed to it.
func (r *Rantai) Exclude(middleware ...func(http.Handler) http.Handler) *Rantai {
	newMiddlewares := make([]func(http.Handler) http.Handler, 0)

	for _, m := range r.Middlewares {
		for _, ex := range middleware {
			if !kmpr.Do(m, ex) {
				newMiddlewares = append(newMiddlewares, m)
			}
		}
	}

	return &Rantai{
		Middlewares: newMiddlewares,
	}
}
