package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Cors holds the functions and data configured and provide the middleware
// used for CORS (Cross-origin resource sharing).
type Cors struct {
	allowedOrigins string
	allowedHeaders string
	allowedMethods string
	maxAge         string
}

// ConfigFunc is the type of function used to configure the Cors
// instance. The library provide various functions that return ConfigFunc
// compatible functions.
type ConfigFunc func(*Cors)

// New creates a new Cors instance that is configured with the given
// ConfigFunc.
func New(configs ...ConfigFunc) *Cors {
	c := &Cors{}

	for _, cFn := range configs {
		cFn(c)
	}

	return c
}

func (c *Cors) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.allowedOrigins != "" {
			w.Header().Add("Access-Control-Allow-Origin", c.allowedOrigins)
		}
		if c.allowedMethods != "" {
			w.Header().Add("Access-Control-Allow-Methods", c.allowedMethods)
		}
		if c.allowedHeaders != "" {
			w.Header().Add("Access-Control-Allow-Headers", c.allowedHeaders)
		}

		if r != nil && r.Method == http.MethodOptions {
			if c.maxAge != "" {
				w.Header().Add("Access-Control-Max-Age", c.maxAge)
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// WithOrigins returns a ConfigFunc that configures the Cors to output a
// header that signals that only requests from the given hosts are accepted.
func WithOrigins(origins ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedOrigins = strings.Join(origins, ", ")
	}
}

// WithMethods returns a ConfigFunc that configures the Cors to output
// a header that signals that only requests with one of the given methods
// are accepted.
func WithMethods(methods ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedMethods = strings.Join(methods, ", ")
	}
}

// WithMaxAge returns a ConfigFunc that configures the Cors to output
// a header that signals that the CORS information (optained from a
// request method OPTIONS) could be cached for the given amount of time.
func WithMaxAge(age time.Duration) ConfigFunc {
	return func(c *Cors) {
		c.maxAge = fmt.Sprint(int(age.Seconds()))
	}
}

// WithHeaders returns a ConfigFunc that configures the Cors to output
// a header that signals that only the given headers are accepted.
func WithHeaders(headers ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedHeaders = strings.Join(headers, ", ")
	}
}
