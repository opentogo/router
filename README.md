# router

[![Build Status](https://travis-ci.org/opentogo/router.svg?branch=master)](https://travis-ci.org/opentogo/router)
[![GoDoc](https://godoc.org/github.com/opentogo/router?status.png)](https://godoc.org/github.com/opentogo/router)
[![codecov](https://codecov.io/gh/opentogo/router/branch/master/graph/badge.svg)](https://codecov.io/gh/opentogo/router)
[![Go Report Card](https://goreportcard.com/badge/github.com/opentogo/router)](https://goreportcard.com/report/github.com/opentogo/router)
[![Open Source Helpers](https://www.codetriage.com/opentogo/router/badges/users.svg)](https://www.codetriage.com/opentogo/router)

A simple HTTP router in Go.

## Installation

```bash
go get github.com/opentogo/router
```

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/opentogo/router"
)

func main() {
	r := &router.Router{}
	r.Handler(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("hello, world!")); err != nil {
			http.Error(w, "Unable to write response for root handler.", http.StatusInternalServerError)
		}
	})

	r.NotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write(nil); err != nil {
			http.Error(w, "Unable to write response for NotFound handler.", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      r,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
```

## Contributors

- [rogeriozambon](https://github.com/rogeriozambon) Rogério Zambon - creator, maintainer
