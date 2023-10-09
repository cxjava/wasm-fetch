package fetch

import (
	"context"
	"io"
)

// Opts are the options you can pass to the fetch call.
type Opts struct {
	// Method is the http verb (constants are copied from net/http to avoid import)
	Method string

	// Headers is a map of http headers to send.
	Headers map[string]string

	// Body is the body request
	Body io.Reader

	// Mode docs https://developer.mozilla.org/en-US/docs/Web/API/Request/mode
	Mode string

	// Credentials docs https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
	Credentials string

	// Cache docs https://developer.mozilla.org/en-US/docs/Web/API/Request/cache
	Cache string

	// Redirect docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	Redirect string

	// Referrer docs https://developer.mozilla.org/en-US/docs/Web/API/Request/referrer
	Referrer string

	// ReferrerPolicy docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	ReferrerPolicy string

	// Integrity docs https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
	Integrity string

	// KeepAlive docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	KeepAlive *bool

	// Signal docs https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal
	Signal context.Context
}
