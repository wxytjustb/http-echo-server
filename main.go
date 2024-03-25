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

type LogData struct {
	ApisixLatency   float64  `json:"apisix_latency"`
	ServiceID       string   `json:"service_id"`
	Server          Server   `json:"server"`
	ClientIP        string   `json:"client_ip"`
	UpstreamLatency float64  `json:"upstream_latency"`
	Latency         float64  `json:"latency"`
	StartTime       int64    `json:"start_time"`
	Response        Response `json:"response"`
	RouteID         string   `json:"route_id"`
	Upstream        string   `json:"upstream"`
	Request         Request  `json:"request"`
}

type Server struct {
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

type Response struct {
	Size    int     `json:"size"`
	Status  int     `json:"status"`
	Body    string  `json:"body"`
	Headers Headers `json:"headers"`
}

type Headers struct {
	ContentLength string `json:"content-length"`
	Server        string `json:"server"`
	Connection    string `json:"connection"`
	ContentType   string `json:"content-type"`
}

type Request struct {
	URL         string                 `json:"url"`
	Querystring map[string]interface{} `json:"querystring"`
	Size        int                    `json:"size"`
	URI         string                 `json:"uri"`
	Method      string                 `json:"method"`
	Headers     Headers                `json:"headers"`
}

func handleLogRequest(c *gin.Context) {

	var logData []LogData
	err := json.NewDecoder(c.Request.Body).Decode(&logData)
	if err != nil {
		// handle error here
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, data := range logData {
		if data.Response.Status >= 400 {
			jsonData, err := json.MarshalIndent(data, "", "    ")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(string(jsonData))
		}
	}

}

func handleDefaultRequests(c *gin.Context) {
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

	c.String(http.StatusOK, string(jsonData)) // Handle all other requests
}

func main() {
	r := gin.Default()

	// Create a new route group
	wildcardGroup := r.Group("/")
	wildcardGroup.Any("/*path", func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/log" {
			// Handle /log POST request
			handleLogRequest(c)
		} else {
			// Handle all other requests
			handleDefaultRequests(c)
		}
	})

	r.Run(":8889")
}
