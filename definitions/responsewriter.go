package definitions

import (
	"fmt"
	"net/http"
)

// New creates a CustomResponseWriter with defaults.
// This uses a receiver to set up a CustomResponseWriter conveniently.
func (CustomResponseWriter) New(w http.ResponseWriter) CustomResponseWriter {
	cw := CustomResponseWriter{
		ResponseWriter: w,
	}
	if flusher, ok := w.(http.Flusher); ok {
		cw.flusher = flusher
	}
	return cw
}

func (cw *CustomResponseWriter) WriteHeader(statusCode int) {
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *CustomResponseWriter) Write(data []byte) (int, error) {
	cw.flusher.Flush()
	return cw.ResponseWriter.Write(data)

}

func (cw *CustomResponseWriter) Flush() {
	if cw.flusher != nil {
		cw.flusher.Flush()
	}
}

func (cw *CustomResponseWriter) Stream(data string) {
	// Set CORS headers
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	cw.ResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	// Set the Content-Type to text/event-stream for streaming updates
	cw.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	cw.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	cw.ResponseWriter.Header().Set("Connection", "keep-alive")

	// Write the data in the SSE (Server-Sent Events) format
	fmt.Fprintf(cw.ResponseWriter, "data: %s\n\n", data)
	if cw.flusher != nil {
		cw.flusher.Flush()
	}
}

// New function to stream a Payload
func (cw *CustomResponseWriter) StreamPayload(payload Payload) {
	// Set CORS headers (if not already set)
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	cw.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	cw.ResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	// Set the Content-Type to text/event-stream for streaming updates
	cw.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	cw.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	cw.ResponseWriter.Header().Set("Connection", "keep-alive")

	// Format the event and message
	fmt.Fprintf(cw.ResponseWriter, "event: %s\ndata: %s\n\n", payload.Event, payload.Message)

	// Use a single `fmt.Fprintf` statement to send the payload as SSE

	// Flush the data immediately
	if cw.flusher != nil {
		cw.flusher.Flush()
	}
}
