package web

import (
	"net/http"
)

type GinRouter struct {
	routes     map[string]map[string]HandlerFunc
	middleware []Middleware
}

type GinContext struct {
	request  *http.Request
	response *GinResponseWriter
	params   map[string]string
	data     map[string]interface{}
}

type GinResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
	written bool
}

func NewGinRouter() *GinRouter {
	return &GinRouter{
		routes: map[string]map[string]HandlerFunc{
			"GET":    make(map[string]HandlerFunc),
			"POST":   make(map[string]HandlerFunc),
			"PUT":    make(map[string]HandlerFunc),
			"DELETE": make(map[string]HandlerFunc),
		},
		middleware: make([]Middleware, 0),
	}
}

func (r *GinRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &GinContext{
		request:  req,
		response: &GinResponseWriter{ResponseWriter: w},
		params:   make(map[string]string),
		data:     make(map[string]interface{}),
	}
	
	handler := r.findHandler(req.Method, req.URL.Path)
	if handler == nil {
		http.NotFound(w, req)
		return
	}
	
	finalHandler := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		finalHandler = r.middleware[i](finalHandler)
	}
	
	if err := finalHandler(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *GinRouter) findHandler(method, path string) HandlerFunc {
	methodRoutes, exists := r.routes[method]
	if !exists {
		return nil
	}
	
	handler, exists := methodRoutes[path]
	if !exists {
		return nil
	}
	
	return handler
}

func (r *GinRouter) GET(path string, handler HandlerFunc) {
	r.routes["GET"][path] = handler
}

func (r *GinRouter) POST(path string, handler HandlerFunc) {
	r.routes["POST"][path] = handler
}

func (r *GinRouter) PUT(path string, handler HandlerFunc) {
	r.routes["PUT"][path] = handler
}

func (r *GinRouter) DELETE(path string, handler HandlerFunc) {
	r.routes["DELETE"][path] = handler
}

func (r *GinRouter) Use(middleware ...Middleware) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *GinRouter) Group(prefix string) RouterGroup {
	return &GinRouterGroup{
		router: r,
		prefix: prefix,
	}
}

type GinRouterGroup struct {
	router *GinRouter
	prefix string
}

func (g *GinRouterGroup) GET(path string, handler HandlerFunc) {
	g.router.GET(g.prefix+path, handler)
}

func (g *GinRouterGroup) POST(path string, handler HandlerFunc) {
	g.router.POST(g.prefix+path, handler)
}

func (g *GinRouterGroup) PUT(path string, handler HandlerFunc) {
	g.router.PUT(g.prefix+path, handler)
}

func (g *GinRouterGroup) DELETE(path string, handler HandlerFunc) {
	g.router.DELETE(g.prefix+path, handler)
}

func (g *GinRouterGroup) Use(middleware ...Middleware) {
	g.router.Use(middleware...)
}

func (g *GinRouterGroup) Group(prefix string) RouterGroup {
	return &GinRouterGroup{
		router: g.router,
		prefix: g.prefix + prefix,
	}
}


func (c *GinContext) Request() *http.Request {
	return c.request
}

func (c *GinContext) Response() ResponseWriter {
	return c.response
}

func (c *GinContext) Param(key string) string {
	return c.params[key]
}

func (c *GinContext) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *GinContext) Bind(obj interface{}) error {
	return nil
}

func (c *GinContext) JSON(code int, obj interface{}) error {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(code)
	
	c.response.Write([]byte(`{"status":"ok"}`))
	return nil
}

func (c *GinContext) String(code int, format string, values ...interface{}) error {
	c.response.Header().Set("Content-Type", "text/plain")
	c.response.WriteHeader(code)
	
	c.response.Write([]byte(format))
	return nil
}

func (c *GinContext) File(filepath string) error {
	http.ServeFile(c.response, c.request, filepath)
	return nil
}

func (c *GinContext) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *GinContext) Get(key string) interface{} {
	return c.data[key]
}


func (w *GinResponseWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *GinResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *GinResponseWriter) Status() int {
	return w.status
}

func (w *GinResponseWriter) Size() int {
	return w.size
}

func (w *GinResponseWriter) Written() bool {
	return w.written
}
