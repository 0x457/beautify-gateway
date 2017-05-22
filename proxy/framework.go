package proxy

import (
	"net/http"
	"sync"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/ivessong/beautify-gateway/models"
)

//Framework  proxy 服务
type Framework struct {
	middleware      Middleware
	apisRouter      *mux.Router
	apisRouterMutex sync.RWMutex
	rt              *models.RouteTable
	contextPool     ContextPool
}

//NewProxy  初始化proxy
func NewProxy(rt *models.RouteTable) *Framework {
	f := &Framework{
		rt:         rt,
		middleware: make([]Handler, 0),
	}
	f.contextPool = &contextPool{
		sync.Pool{New: func() interface{} { return &Context{framework: f} }},
	}
	f.init()
	return f
}

//Use  过滤器
func (f *Framework) Use(handles ...HandlerFunc) {
	f.middleware = joinMiddleware(f.middleware, convertToHandlers(handles))
}

func (f *Framework) init() {
	f.apisRouter = mux.NewRouter()
	url := "/products/{key}"
	f.apisRouter.NewRoute().Name(url).Path(url)
}

//reloadapisRouter  重新加载route
func (f *Framework) reloadapisRouter() {
}

var wgPool = &sync.Pool{New: func() interface{} {
	return &sync.WaitGroup{}
}}

//Serve  server http
func (f *Framework) Serve(w http.ResponseWriter, r *http.Request) {
	var match mux.RouteMatch
	if ok := f.match(r, &match); !ok {
		//TODO: 无路由匹配
		fmt.Printf("there is no mux route for this \n")
		return
	}
	fmt.Println(match.Route.GetName(), match.Vars)
	var routeResults []*models.RouteResult
	if routeResults = f.rt.Get(r, match.Route.GetName()); routeResults == nil {
		//TODO:未找到相应的api
		fmt.Printf("there is no api route for this \n")
		return
	}
	//开始处理middleware
	ctx, _ := f.contextPool.Acquire(w, r, match.Vars)
	ctx.RouteResults = routeResults
	defer f.contextPool.Release(ctx)
	ctx.Do()
	if ctx.ErrEnd { //中间件验证是否结束
		return
	}
	count := len(routeResults)
	merge := count > 1
	if merge {
		wg := wgPool.Get().(*sync.WaitGroup)
		wg.Add(count)
		for _, item := range routeResults {
			item.Merge = merge
			go func(item *models.RouteResult) {
				f.handle(ctx, wg, item)
			}(item)
		}
		wg.Wait()
	} else {
		f.handle(ctx, nil, routeResults[0])
	}
	//开始处理以及组合数据
}

//Match  handle match proxy request
func (f *Framework) match(r *http.Request, match *mux.RouteMatch) bool {
	defer f.apisRouterMutex.RUnlock()
	f.apisRouterMutex.RLock()
	if ok := f.apisRouter.Match(r, match); !ok { //路由匹配
		return false
	}
	return true
}

func (f *Framework) handle(ctx *Context, wg *sync.WaitGroup, route *models.RouteResult) {

}
