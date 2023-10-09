package fetch

// Response is the response that return from the fetch promise.
type Response struct {
	Headers    Header
	OK         bool
	Redirected bool
	Status     int
	StatusText string
	Type       string
	URL        string
	Body       []byte
	BodyUsed   bool
}
