package ctxrouter

import (
	"net/http"

	"github.com/corestoreio/csfw/net/mw"
)

// Group represents sub router bound to a path prefix
type Group struct {
	r Router
}

func (g *Group) Use(mws ...mw.Middleware) {
	g.r.middleware = append(g.r.middleware, mws...)
}

//func (g *Group) CONNECT(path string, h http.Handler) {
//	g.r.CONNECT(path, h)
//}

func (g *Group) DELETE(path string, h http.HandlerFunc) {
	g.r.DELETE(path, h)
}

func (g *Group) GET(path string, h http.HandlerFunc) {
	g.r.GET(path, h)
}

func (g *Group) HEAD(path string, h http.HandlerFunc) {
	g.r.HEAD(path, h)
}

func (g *Group) OPTIONS(path string, h http.HandlerFunc) {
	g.r.OPTIONS(path, h)
}

func (g *Group) PATCH(path string, h http.HandlerFunc) {
	g.r.PATCH(path, h)
}

func (g *Group) POST(path string, h http.HandlerFunc) {
	g.r.POST(path, h)
}

func (g *Group) PUT(path string, h http.HandlerFunc) {
	g.r.PUT(path, h)
}

func (g *Group) WEBSOCKET(path string, h http.HandlerFunc) {
	g.r.WEBSOCKET(path, h)
}

func (g *Group) ServeFiles(path string, root http.FileSystem) {
	g.r.ServeFiles(path, root)
}

func (g *Group) Group(prefix string, mws ...mw.Middleware) *Group {
	return g.r.Group(prefix, mws...)
}
