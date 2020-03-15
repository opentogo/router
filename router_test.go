package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allisson/go-assert"
)

func TestRouter(t *testing.T) {
	router := &Router{}

	t.Run("HTTP methods", func(t *testing.T) {
		cases := []struct {
			description    string
			method         string
			expectedMethod int
		}{
			{
				description:    "GET method",
				method:         http.MethodGet,
				expectedMethod: http.StatusOK,
			},
			{
				description:    "POST method",
				method:         http.MethodPost,
				expectedMethod: http.StatusOK,
			},
			{
				description:    "DELETE method",
				method:         http.MethodDelete,
				expectedMethod: http.StatusOK,
			},
			{
				description:    "PATCH method",
				method:         http.MethodPatch,
				expectedMethod: http.StatusOK,
			},
			{
				description:    "PUT method",
				method:         http.MethodPut,
				expectedMethod: http.StatusOK,
			},
			{
				description:    "HEAD method",
				method:         http.MethodHead,
				expectedMethod: http.StatusOK,
			},
		}
		for _, test := range cases {
			t.Run(test.description, func(t *testing.T) {
				var (
					w = httptest.NewRecorder()
					r = httptest.NewRequest(test.method, "/togo-router", nil)
				)

				router.Handler(test.method, "/togo-router", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "hello, togo/router") // nolint: errcheck
				})
				router.ServeHTTP(w, r)

				assert.Equal(t, test.expectedMethod, w.Code)
				assert.Equal(t, "hello, togo/router", w.Body.String())
			})
		}
	})

	t.Run("not found method", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/togo-super-router", nil)
		)

		router.Handler(http.MethodPost, "/togo-super-router", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "togo-router") // nolint: errcheck
		})
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "404 page not found\n", w.Body.String())
	})

	t.Run("root path", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		router.Handler(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "togo-router") // nolint: errcheck
		})
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "togo-router", w.Body.String())
	})

	t.Run("calling custom `NotFound` handler", func(t *testing.T) {
		var (
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/togo-invalid-router", nil)
		)

		router.NotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found") // nolint: errcheck
		})
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "not found", w.Body.String())
	})

	t.Run("checking url parameters", func(t *testing.T) {
		var (
			captured map[string]string
			w        = httptest.NewRecorder()
			r        = httptest.NewRequest(http.MethodGet, "/togo/415/users/9v02", nil)
		)

		router.Handler(http.MethodGet, "/togo/:id/users/:user_id", func(w http.ResponseWriter, r *http.Request) {
			captured = Params(r)
		})
		router.ServeHTTP(w, r)

		assert.Equal(t, "415", captured["id"])
		assert.Equal(t, "9v02", captured["user_id"])
		assert.Equal(t, "", captured["invalid"])
	})

	t.Run("checking url parameters without named parameters", func(t *testing.T) {
		var (
			captured map[string]string
			w        = httptest.NewRecorder()
			r        = httptest.NewRequest(http.MethodGet, "/togo/router", nil)
		)

		router.Handler(http.MethodGet, "/togo/router", func(w http.ResponseWriter, r *http.Request) {
			captured = Params(r)
		})
		router.ServeHTTP(w, r)

		assert.Nil(t, captured)
	})
}
