package main

type ResponseBuffer struct {
	// Embed bytes.Buffer in
	bytes.Buffer

	header map[string][]string
	code   int
}

func NewResponseBuffer() *ResponseBuffer {
	return &ResponseBuffer{
		header: map[string][]string{},
		code:   200, // default OK
	}
}

func (b *ResponseBuffer) Header() http.Header {
	return b.header
}

func (b *ResponseBuffer) WriteHeader(code int) {
	b.code = code
}

func (b *ResponseBuffer) Flush(w http.ResponseWriter) {
	// Grab the header, which is a map pointer
	header := w.Header()

	// Copy the dummy header values to the actual header
	for k, v := range b.header {
		header[k] = v
	}

	// Write the header
	w.WriteHeader(b.code)

	// Copy the body
	io.Copy(w, b)
}