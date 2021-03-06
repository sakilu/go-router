package main

import (
	"fmt"
	"github.com/sakilu/go-router"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"wizard/controllers"
)

func main() {
	// 初始化路由器
	var myRouter router.ControllerRegistor

	// 註冊網頁靜態頁面路徑
	myRouter.SetStaticPath("/", "/html/")

	// 路由器 controller 註冊開始
	myRouter.Add("/", &controllers.IndexController{})

	// 註冊結束 啟動伺服器
	http.HandleFunc("/", myRouter.ServeHTTP)

	open.Run("http://127.0.0.1:8123/index.html")
	err := http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Println(err)
	}
}
