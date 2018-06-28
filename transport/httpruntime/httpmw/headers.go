/*Package httpmw provides middlewares that are automatically
used by the generated code.*/
package httpmw

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

// DefaultChain is a chain that gets applied to the generated handlers.
func DefaultChain(next http.HandlerFunc) http.HandlerFunc {
	return HeadersToGRPCMD(next)
}

// HeadersToGRPCMD inserts HTTP headers to gRPC metadata, as if they were
// received via gRPC.
// Every header name is lowercased, per gRPC standards.
func HeadersToGRPCMD(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// use metadata.FromIncomingContext to access it
		var md metadata.MD

		ctx := r.Context()
		// Use existing MD if it was injected earlier
		if m, ok := metadata.FromIncomingContext(ctx); ok {
			md = m
		} else {
			md = make(metadata.MD)
		}

		for k, v := range r.Header {
			md.Append(k, v...)
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
