package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Any("/*path", func(c *gin.Context) {
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			headers[k] = v[0]
		}

		data := gin.H{
			"method":    c.Request.Method,
			"uri":       c.Request.RequestURI,
			"body":      c.Request.Body,
			"headers":   headers,
			"ip":        c.ClientIP(),
			"remote_ip": c.RemoteIP(),
			"host":      c.Request.Host,
		}

		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(jsonData))

		c.String(http.StatusOK, string(jsonData))

		//c.JSON(http.StatusOK, data)
	})

	r.Run(":8889")
}
