package router

import (
	"context"
	"net/http"
	"regexp"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type Router struct {
	routes []route
}

// regexp example -> https://regex101.com/r/84S9iL/1
func (router *Router) NewRoute(method, regexpString string, handler http.HandlerFunc) {
	regex := regexp.MustCompile("^" + regexpString + "$")

	router.routes = append(router.routes, route{
		method,
		regex,
		handler,
	})
}

func (router *Router) Serve(w http.ResponseWriter, r *http.Request) {
	for _, v := range router.routes {
		mmaatchtch := v.regex.FindStringSubmatch(r.URL.Path)

		if len(mmaatchtch) > 0 {
			if r.Method != v.method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			matchMap := make(map[string]string)
			groupName := v.regex.SubexpNames()

			// map group name(key) to submatched result
			// these arrays have one to one relationship
			for i := 1; i < len(mmaatchtch); i++ {
				matchMap[groupName[i]] = mmaatchtch[i]
			}

			ctx := context.WithValue(r.Context(), struct{}{}, matchMap)
			v.handler(w, r.WithContext(ctx))
			return
		}
	}
}

func GetField(r *http.Request, name string) string {
	fields := r.Context().Value(struct{}{}).(map[string]string)
	return fields[name]
}
