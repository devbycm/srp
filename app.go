package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	new(App).
		Action("proxy :public :tunnel", proxy).
		Action("server :tunnel :service", server).
		Run()
}

func proxy(c *Context) {
	p := NewProxy(c.Get("public"), c.Get("tunnel"))
	check(p.Run())
}

func server(c *Context) {
	tunnelAddr := c.Get("tunnel")
	serviceAddr := c.Get("service")

	s := NewServer(tunnelAddr, serviceAddr)
	check(s.CreateService(), "connect.service")
	check(s.ConnectTunnel(), "connect.tunnel")

	r := gin.New()
	r.GET("", func(c *gin.Context) {
		data := c.Query("x")
		c.String(http.StatusOK, data)
		fmt.Println(data)
	})
	check(http.Serve(s.Listener, r))
}
