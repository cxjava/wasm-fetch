// Package fetch is a  Web Assembly fetch wrapper that avoids importing net/http.
/*
	package main

	import (
		"context"
		"time"

		"github.com/cxjava/wasm-fetch"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := fetch.Fetch("/some/api/call", &fetch.Opts{
		Body:   strings.NewReader(`{"one": "two"}`),
		Method: fetch.MethodPost,
		Signal: ctx,
	})
*/
package fetch // import "github.com/cxjava/wasm-fetch"
