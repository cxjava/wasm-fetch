//go:build wasm
package fetch

import (
	"bytes"
	"errors"
	"io"
	"syscall/js"
)

// Fetch uses the JS Fetch API to make requests
// over WASM.
func Fetch(url string, opts *Opts) (*Response, error) {
	optsMap, err := mapOpts(opts)
	if err != nil {
		return nil, err
	}

	type fetchResponse struct {
		r *Response
		b js.Value
		e error
	}
	ch := make(chan *fetchResponse)
	done := make(chan struct{}, 1)
	if opts.Signal != nil {
		controller := js.Global().Get("AbortController").New()
		signal := controller.Get("signal")
		optsMap["signal"] = signal
		go func() {
			select {
			case <-opts.Signal.Done():
				controller.Call("abort")
			case <-done:
			}
		}()
	}

	success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		r := new(Response)
		resp := args[0]
		headersIt := resp.Get("headers").Call("entries")
		headers := Header{}
		for {
			n := headersIt.Call("next")
			if n.Get("done").Bool() {
				break
			}
			pair := n.Get("value")
			key, value := pair.Index(0).String(), pair.Index(1).String()
			headers.Add(key, value)
		}
		r.Headers = headers
		r.OK = resp.Get("ok").Bool()
		r.Redirected = resp.Get("redirected").Bool()
		r.Status = resp.Get("status").Int()
		r.StatusText = resp.Get("statusText").String()
		r.Type = resp.Get("type").String()
		r.URL = resp.Get("url").String()
		r.BodyUsed = resp.Get("bodyUsed").Bool()

		ch <- &fetchResponse{r: r, b: resp}
		return nil
	})
	defer success.Release()

	failure := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msg := args[0].Get("message").String()
		done <- struct{}{}
		ch <- &fetchResponse{e: errors.New(msg)}
		return nil
	})
	defer failure.Release()

	go js.Global().Call("fetch", url, optsMap).Call("then", success).Call("catch", failure)

	r := <-ch
	if r.e != nil {
		return nil, r.e
	}

	successBody := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Wrap the input ArrayBuffer with a Uint8Array
		uint8arrayWrapper := js.Global().Get("Uint8Array").New(args[0])
		r.r.Body = make([]byte, uint8arrayWrapper.Get("byteLength").Int())
		js.CopyBytesToGo(r.r.Body, uint8arrayWrapper)
		ch <- r
		return nil
	})
	defer successBody.Release()

	failureBody := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Assumes it's a TypeError. See
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/TypeError
		// for more information on this type.
		// See https://fetch.spec.whatwg.org/#concept-body-consume-body for error causes.
		msg := args[0].Get("message").String()
		ch <- &fetchResponse{e: errors.New(msg)}
		return nil
	})
	defer failureBody.Release()

	go r.b.Call("arrayBuffer").Call("then", successBody, failureBody)

	r = <-ch
	return r.r, r.e
}

// oof.
func mapOpts(opts *Opts) (map[string]interface{}, error) {
	mp := map[string]interface{}{}

	if opts.Method != "" {
		mp["method"] = opts.Method
	}
	if opts.Headers != nil {
		mp["headers"] = mapHeaders(opts.Headers)
	}
	if opts.Mode != "" {
		mp["mode"] = opts.Mode
	}
	if opts.Credentials != "" {
		mp["credentials"] = opts.Credentials
	}
	if opts.Cache != "" {
		mp["cache"] = opts.Cache
	}
	if opts.Redirect != "" {
		mp["redirect"] = opts.Redirect
	}
	if opts.Referrer != "" {
		mp["referrer"] = opts.Referrer
	}
	if opts.ReferrerPolicy != "" {
		mp["referrerPolicy"] = opts.ReferrerPolicy
	}
	if opts.Integrity != "" {
		mp["integrity"] = opts.Integrity
	}
	if opts.KeepAlive != nil {
		mp["keepalive"] = *opts.KeepAlive
	}

	if opts.Body != nil {
		var bts []byte
		var err error

		switch v := opts.Body.(type) {
		case *bytes.Buffer:
			bts = v.Bytes()
		default:
			bts, err = io.ReadAll(opts.Body)
		}

		if err != nil {
			return nil, err
		}

		jsbody := js.Global().Get("Uint8Array").New(len(bts))
		js.CopyBytesToJS(jsbody, bts)

		mp["body"] = jsbody
	}

	return mp, nil
}

func mapHeaders(mp map[string]string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range mp {
		newMap[k] = v
	}
	return newMap
}
