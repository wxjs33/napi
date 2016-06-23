package router

import (
	"net/http"
	"github.com/wxjs33/napi/log"
)

type Router struct {
	handlers map[string]http.Handler
	log *log.Log
}

func InitRouter(log *log.Log) *Router {
	r := &Router{}
	r.handlers = make(map[string]http.Handler)
	r.log = log

	return r
}

func (r *Router) AddRouter(url string, handler http.Handler) error {
	if _, ok := r.handlers[url]; ok {
		r.log.Error("url: %s has been added", url)
		//TODO: add some error
		return nil
	}
	r.handlers[url] = handler

	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h, ok := r.handlers[req.URL.Path]; ok {
		h.ServeHTTP(w, req)
		//http.Error(w, "Method Not Allowed", 405)

	} else {
		//if r.NotFound != nil {
		//	r.NotFound.ServeHTTP(w, req)
		//	return
		//}

		//logger.Info.Printf("%s Not Found", req.URL.Path)
		http.Error(w, "URL Not Found", 404)
	}
}


