//go:build !wasm
package fetch

func Fetch(url string, opts *Opts) (*Response, error) {
	panic("not to be used")
}
