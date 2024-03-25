package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	r := gin.Default()

	r.Any("/*path", func(c *gin.Context) {
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			headers[k] = v[0]
		}

		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// After reading, we have to set it back
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		bodyString := string(bodyBytes)

		uri, _ := url.QueryUnescape(c.Request.RequestURI)
		data := gin.H{
			"method":    c.Request.Method,
			"uri":       uri,
			"body":      bodyString,
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
