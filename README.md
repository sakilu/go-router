go-router
=========
<p>
修改自 beego router. <br />
出處 <br />
https://github.com/astaxie/beego <br />
</p>


<code>
	// 初始化路由器
	var myRouter router.ControllerRegistor
	// 註冊網頁靜態頁面路徑
	myRouter.SetStaticPath("/", "/html/")

	// 路由器 controller 註冊開始
	myRouter.Add("/", &controllers.IndexController{})

	// 註冊結束 啟動伺服器
	http.HandleFunc("/", myRouter.ServeHTTP)

	open.Run("http://127.0.0.1:8123")
	http.ListenAndServe(":8123", nil)

</code>
