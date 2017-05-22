package models

import (
	"net/http"
	"sync"
)

//RouteTable  路由信息表
//根据 mux route
type RouteTable struct {
	rwLock      *sync.RWMutex //读写锁
	routes      map[string]*Route
	requestPool *sync.Pool
}

//newRouteTable  new a route
func newRouteTable() *RouteTable {
	return &RouteTable{
		requestPool: &sync.Pool{New: func() interface{} { return &http.Request{} }},
	}
}

//ReleaseRoute  release route request
func (r *RouteTable) ReleaseRoute(result *RouteResult) {
	r.requestPool.Put(result.Request)
}

//Get route
func (r *RouteTable) Get(request *http.Request, path string) []*RouteResult {
	return nil
}

//Update routetable
func (r *RouteTable) Update() {

}

//Add routetable
func (r *RouteTable) Add() {

}

//Del routetable
func (r *RouteTable) Del() {

}
