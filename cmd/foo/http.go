package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"experiments/goa-cors-override/gen/foo"
	"experiments/goa-cors-override/gen/http/foo/server"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, fooEndpoints *foo.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		fooServer *server.Server
	)
	{
		eh := errorHandler(logger)
		fooServer = server.New(fooEndpoints, mux, dec, enc, eh, nil)
		if debug {
			servers := goahttp.Servers{
				fooServer,
			}
			servers.Use(httpmdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	// foosvr.Mount(mux, fooServer)
	// We can't call the foosvr.Mount as is since it will attempt to mount both the CORS handler
	// as well as the FooOptions handler and that will cause a panic at runtime. We mount a custom
	// handler that handles all options requests. This custom handler will route the corresponding
	// options request to our handler and send the rest to the CORS handler where the preflight
	// can happen.
	server.MountFoo1Handler(mux, fooServer.Foo1)
	server.MountFoo2Handler(mux, fooServer.Foo2)
	server.MountFoo3Handler(mux, fooServer.Foo3)
	mux.Handle("OPTIONS", "/*", customOptionsHandler(server.HandleFooOrigin(fooServer.CORS), fooServer.FooOptions))

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpmdlwr.Log(adapter)(handler)
		handler = httpmdlwr.RequestID()(handler)
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler}
	for _, m := range fooServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Printf("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}

// customOptionsHandler routes the OPTIONS request based on the presence of the origin header.
func customOptionsHandler(ch http.Handler, fh http.Handler) http.HandlerFunc {
	corsHandler := ch.(http.HandlerFunc)
	 fooHandler := fh.(http.HandlerFunc)
	return func(w http.ResponseWriter, r *http.Request) {
		// As an example, here we route it based on the presence of the Origin header. It can even
		// be based on the request path or other request properties.
		if r.Method == http.MethodOptions && r.Header.Get("Origin") == "" {
			fmt.Printf("Routing to custom CORS handler because no Origin header is present\n")
			fooHandler(w, r)
			return
		}
		fmt.Printf("Handling it as a CORS preflight request\n")
		fmt.Println(r)
		corsHandler(w, r)
		return
	}
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}
