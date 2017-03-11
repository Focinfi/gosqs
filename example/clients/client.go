package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	handler := func(c *gin.Context) {

		fmt.Printf("%#v\n", c.PostForm("message"))
		c.String(http.StatusOK, "")
	}

	server.POST("/greeting/1", handler)
	server.POST("/greeting/2", handler)
	server.POST("/greeting/3", handler)
	server.Run(":55446")
}
