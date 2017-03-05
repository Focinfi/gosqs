package main

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	handler := func(c *gin.Context) {
		param := &struct {
			Message string `json:"message"`
		}{}

		c.BindJSON(param)
		fmt.Printf("%#v", param)
		c.String(http.StatusOK, "")
	}

	server.POST("/greeting/1", handler)
	server.POST("/greeting/2", handler)
	server.POST("/greeting/3", handler)
	server.Run(":55446")
}
