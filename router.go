package router

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Controller struct {
}

type ControllerInterface interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) Put(w http.ResponseWriter, r *http.Request) {

}

func (cr *ControllerRegistor) SetStaticPath(url string, path string) *ControllerRegistor {
	if cr.StaticDir == nil {
		cr.StaticDir = make(map[string]string)
	}
	cr.StaticDir[url] = path
	return cr
}

type controllerInfo struct {
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
}

type ControllerRegistor struct {
	StaticDir map[string]string
	routers   []*controllerInfo
}

func (p *ControllerRegistor) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"

			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}

	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {

		panic(regexErr)
		return
	}

	t := reflect.Indirect(reflect.ValueOf(c)).Type()
	route := &controllerInfo{}
	route.regex = regex
	route.params = params
	route.controllerType = t

	p.routers = append(p.routers, route)

}

func (p *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var started bool
	for prefix, staticDir := range p.StaticDir {

		if strings.HasPrefix(r.URL.Path, prefix) {
			file := "." + staticDir + r.URL.Path[len(prefix):]
			finfo, err := os.Stat(file)
			if err != nil || finfo.IsDir() {
				continue
			}
			if strings.HasSuffix(file, ".html") || strings.HasSuffix(file, ".htm") {
				w.Header().Set("Content-Type", "text/html")
			}
			http.ServeFile(w, r, file)
			started = true
			return
		}
	}
	requestPath := r.URL.Path

	for _, route := range p.routers {

		if !route.regex.MatchString(requestPath) {
			continue
		}

		matches := route.regex.FindStringSubmatch(requestPath)

		if len(matches[0]) != len(requestPath) {
			continue
		}

		params := make(map[string]string)
		if len(route.params) > 0 {
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.params[i], match)
				params[route.params[i]] = match
			}
		}

		vc := reflect.New(route.controllerType)
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(w)
		in[1] = reflect.ValueOf(r)

		if r.Method == "GET" {
			method := vc.MethodByName("Get")
			method.Call(in)
		} else if r.Method == "POST" {
			method := vc.MethodByName("Post")
			method.Call(in)
		} else if r.Method == "DELETE" {
			method := vc.MethodByName("Delete")
			method.Call(in)
		} else if r.Method == "PUT" {
			method := vc.MethodByName("Put")
			method.Call(in)
		}
		started = true
		break
	}

	if started == false {
		http.NotFound(w, r)
	}
}
