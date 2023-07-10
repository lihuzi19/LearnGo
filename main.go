package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.GET(getHelloWorld())
	r.POST(postUploadFile())
	err := r.Run(":5560")
	if err != nil {
		return
	}
}

func getHelloWorld() (string, func(c *gin.Context)) {
	return "hello", responseSuccess("Hello World")
}

func postUploadFile() (string, func(c *gin.Context)) {
	return "uploadFile", func(c *gin.Context) {
		//FormFile返回所提供的表单键的第一个文件
		f, _ := c.FormFile("file")
		var path = "./"
		androidDir, _ := os.Create("/data/data/com.termux/files/usr/share/nginx/html/")
		_, err := androidDir.Stat()
		if err == nil {
			path = "/data/data/com.termux/files/usr/share/nginx/html/"
		}
		err = c.SaveUploadedFile(f, path+f.Filename)
		if err != nil {
			return
		}
		responseSuccessWithContext(c, "Upload File Success")
	}
}

func responseSuccess(data string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": "0",
			"msg":  "OK",
			"data": data,
		})
	}
}
func responseSuccessWithContext(c *gin.Context, data string) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
		"data": data,
	})
}
