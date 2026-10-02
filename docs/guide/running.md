# Running Your Application

## Run

The simplest way to serve an application is its `Run` method:

```go
app := vox.New()
if err := app.Run("localhost:3000"); err != nil {
	log.Fatal(err)
}
```

`Run` takes an address in the form accepted by [`http.ListenAndServe`](https://pkg.go.dev/net/http#ListenAndServe) and blocks until the server fails. Use `":3000"` to listen on all interfaces.

## Using your own http.Server

`vox.Application` implements [`http.Handler`](https://pkg.go.dev/net/http#Handler). For anything beyond the defaults, create the server yourself and hand it the application. Production servers should at least set timeouts:

```go
app := vox.New()

server := &http.Server{
	Addr:              ":3000",
	Handler:           app,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       30 * time.Second,
	WriteTimeout:      30 * time.Second,
	IdleTimeout:       2 * time.Minute,
}
log.Fatal(server.ListenAndServe())
```

### HTTPS

```go
log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
```

### Graceful shutdown

With your own server you can stop accepting connections on a signal and let requests in flight finish:

```go
func main() {
	app := vox.New()
	server := &http.Server{Addr: ":3000", Handler: app}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
```

## Integrating with existing net/http code

Because an application is an `http.Handler`, it can live next to other handlers. This helps when migrating to or from Vox one endpoint at a time.

```go
func rawHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello from a raw handler")
}

func voxHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = "Hello from a vox handler"
}

func main() {
	app := vox.New()
	app.Get("/vox", voxHandler)

	http.HandleFunc("/raw", rawHandler)
	http.Handle("/", app)
	http.ListenAndServe("localhost:3000", nil)
}
```

### Mounting under a prefix

Routes are matched against the request path as the application receives it. To serve an application below a prefix, strip the prefix first:

```go
api := vox.New()
api.Get("/users", listUsers) // served as /api/users

mux := http.NewServeMux()
mux.Handle("/api/", http.StripPrefix("/api", api))
http.ListenAndServe("localhost:3000", mux)
```

### Calling an http.Handler from Vox

The other direction works too. See [Writing the response yourself](./response.md#writing-the-response-yourself).

## Behind a reverse proxy

Vox reports the remote address of the TCP connection, both in `req.RemoteAddr` and in the access log. Behind a proxy or load balancer that is the address of the proxy. Read the forwarding header your proxy sets, such as `X-Forwarded-For`, if you need the address of the client, and only trust it when the request really came through the proxy.

## Profiling

To expose runtime profiling endpoints, see [Pprof](../middlewares/pprof.md).
