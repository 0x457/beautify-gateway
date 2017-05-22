package models

import "net/http"

// api 类型
const (
	ProtocolHTTP = iota
	ProtocolHTTPS
	ProtocolHTTPAndHTTPS
)

//route status
const (
	RouteStatusDown = iota
	RouteStatusUp
)

//Route  api 服务
type Route struct {
	ID          int64
	RateID      int64 //流控
	Name        string
	Description string
	Protocol    int
	Path        string
	Method      string
	Domain      string
	Status      int
	APIs        []*API
}

//RouteResult  route 选举结果
type RouteResult struct {
	API     *API
	Merge   bool          `xorm:"-"`
	Request *http.Request `xorm:"-"`
}
