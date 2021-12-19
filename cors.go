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
	allowedOrigins []string
	allowedHeaders []string
	allowedMethods []string
	maxAge         time.Duration
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
		if len(c.allowedOrigins) > 0 {
			w.Header().Add("Access-Control-Allow-Origin", strings.Join(c.allowedOrigins, ", "))
		}
		if len(c.allowedMethods) > 0 {
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(c.allowedMethods, ", "))
		}
		if len(c.allowedHeaders) > 0 {
			w.Header().Add("Access-Control-Allow-Headers", strings.Join(c.allowedHeaders, ", "))
		}

		if r != nil && r.Method == http.MethodOptions {
			if c.maxAge > 0 {
				w.Header().Add("Access-Control-Max-Age", fmt.Sprint(int(c.maxAge.Seconds())))
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
		c.allowedOrigins = origins
	}
}

// WithMethods returns a ConfigFunc that configures the Cors to output
// a header that signals that only requests with one of the given methods
// are accepted.
func WithMethods(methods ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedMethods = methods
	}
}

// WithMaxAge returns a ConfigFunc that configures the Cors to output
// a header that signals that the CORS information (optained from a
// request method OPTIONS) could be cached for the given amount of time.
func WithMaxAge(age time.Duration) ConfigFunc {
	return func(c *Cors) {
		c.maxAge = age
	}
}

// WithHeaders returns a ConfigFunc that configures the Cors to output
// a header that signals that only the given headers are accepted.
func WithHeaders(headers ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedHeaders = headers
	}
}
