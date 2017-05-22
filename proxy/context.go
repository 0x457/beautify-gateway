package proxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/ivessong/beautify-gateway/models"
)

type (
	// Handler the main Iris Handler interface.
	Handler interface {
		Serve(ctx *Context) // iris-specific
	}
	// HandlerFunc type is an adapter to allow the use of
	// ordinary functions as HTTP handlers.  If f is a function
	// with the appropriate signature, HandlerFunc(f) is a
	// Handler that calls f.
	HandlerFunc func(ctx *Context)

	// Middleware is just a slice of Handler []func(c *Context)
	Middleware []Handler
)

// Serve implements the Handler
func (h HandlerFunc) Serve(ctx *Context) {
	h(ctx)
}

// convertToHandlers just make []HandlerFunc to []Handler, although HandlerFunc and Handler are the same
// we need this on some cases we explicit want a interface Handler, it is useless for users.
func convertToHandlers(handlersFn []HandlerFunc) Middleware {
	hlen := len(handlersFn)
	mlist := make([]Handler, hlen)
	for i := 0; i < hlen; i++ {
		mlist[i] = Handler(handlersFn[i])
	}
	return mlist
}

// joinMiddleware uses to create a copy of all middleware and return them in order to use inside the node
func joinMiddleware(middleware1 Middleware, middleware2 Middleware) Middleware {
	nowLen := len(middleware1)
	totalLen := nowLen + len(middleware2)
	// create a new slice of middleware in order to store all handlers, the already handlers(middleware) and the new
	newMiddleware := make(Middleware, totalLen)
	//copy the already middleware to the just created
	copy(newMiddleware, middleware1)
	//start from there we finish, and store the new middleware too
	copy(newMiddleware[nowLen:], middleware2)
	return newMiddleware
}

type (
	// ContextPool is a set of temporary *Context that may be individually saved and
	// retrieved.
	//
	// Any item stored in the Pool may be removed automatically at any time without
	// notification. If the Pool holds the only reference when this happens, the
	// item might be deallocated.
	//
	// The ContextPool is safe for use by multiple goroutines simultaneously.
	//
	// ContextPool's purpose is to cache allocated but unused Contexts for later reuse,
	// relieving pressure on the garbage collector.
	ContextPool interface {
		// Acquire returns a Context from pool.
		// See Release.
		Acquire(w http.ResponseWriter, r *http.Request, vars map[string]string) (*Context, error)
		// Release puts a Context back to its pull, this function releases its resources.
		// See Acquire.
		Release(ctx *Context)
		// Framework is never used, except when you're in a place where you don't have access to the *iris.Framework station
		// but you need to fire a func or check its Config.
		//
		// Used mostly inside external routers to take the .Config.VHost
		// without the need of other param receivers and refactors when changes
		//
		// note: we could make a variable inside contextPool which would be received by newContextPool
		// but really doesn't need, we just need to borrow a context: we are in pre-build state
		// so the server is not actually running yet, no runtime performance cost.
		Framework() *Framework
	}
	//contextPool  context pool
	contextPool struct {
		pool sync.Pool
	}
)

var _ ContextPool = &contextPool{}

func (c *contextPool) Acquire(w http.ResponseWriter, r *http.Request, vars map[string]string) (*Context, error) {
	ctx := c.pool.Get().(*Context)
	ctx.ResponseWriter = acquireResponseWriter(w)
	ctx.Request = r
	ctx.Vars = vars
	ctx.Host = r.Host
	ctx.Method = r.Method
	ctx.Path = r.URL.Path
	ctx.RawQuery = r.URL.RawQuery
	ctx.RemoteAddress = r.RemoteAddr
	ctx.ContentLength = r.ContentLength
	ctx.Headers = desliceValues(r.Header)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	ctx.Body = string(body)
	if err = r.ParseForm(); err != nil {
		return nil, err
	}
	ctx.Form = desliceValues(r.PostForm)
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	ctx.Query = desliceValues(query)
	return ctx, nil
}

func (c *contextPool) Release(ctx *Context) {
	//TODO:
	ctx.ResponseWriter.releaseMe()
	ctx.ErrEnd = false
	c.pool.Put(ctx)
}

func (c *contextPool) Framework() *Framework {
	ctx := c.pool.Get().(*Context)
	s := ctx.framework
	c.pool.Put(ctx)
	return s
}

//Context  proxy context
type Context struct {
	//context 开始时间 结束时间注入
	startAt  int64
	endAt    int64
	Method   string
	Host     string
	Path     string
	RawQuery string
	Body     string

	RemoteAddress string
	ContentLength int64

	ID string

	Route *models.Route

	RouteResults []*models.RouteResult

	framework *Framework
	// Pos is the position number of the Context, look .Next to understand
	Pos            int // exported because is useful for debugging
	Request        *http.Request
	ResponseWriter *responseWriter

	Vars    map[string]string
	Query   map[string]interface{}
	Headers map[string]interface{}
	Form    map[string]interface{}

	ErrEnd bool //是否结束
}

// Do calls the first handler only, it's like Next with negative pos, used only on Router&MemoryRouter
func (ctx *Context) Do() {
	ctx.Pos = 0
	ctx.framework.middleware[0].Serve(ctx)
}

// Next calls all the next handler from the middleware stack, it used inside a middleware
func (ctx *Context) Next() {
	//set position to the next
	ctx.Pos++
	//run the next
	if ctx.Pos < len(ctx.framework.middleware) {
		ctx.framework.middleware[ctx.Pos].Serve(ctx)
	}
}
