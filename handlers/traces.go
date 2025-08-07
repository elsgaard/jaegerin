package handlers

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
)

func HandleTraces(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/traces", func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/x-protobuf" {
			http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		body, err := io.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}

		var req collectorpb.ExportTraceServiceRequest

		if err := proto.Unmarshal(body, &req); err != nil {
			http.Error(w, "Failed to parse Protobuf", http.StatusBadRequest)
			return
		}

		// Log received spans
		for _, resourceSpan := range req.ResourceSpans {
			for _, scopeSpan := range resourceSpan.ScopeSpans {
				for _, span := range scopeSpan.Spans {
					fmt.Printf("Received span: %s [%d]\n", span.Name, binary.BigEndian.Uint64(span.TraceId))
				}
			}
		}

		// Return an empty successful response
		w.WriteHeader(http.StatusOK)
	})
}
