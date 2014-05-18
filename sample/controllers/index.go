package controllers

import (
	"github.com/sakilu/go-router"
	"net/http"
)

type IndexController struct {
	router.Controller
}

func (c *IndexController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./html/index.html")
}
