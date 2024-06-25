package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.GET(getHelloWorld())
	r.POST(postUploadFile())
	r.GET(apks())
	r.GET(fileList())
	r.StaticFS("/share", http.Dir("D:\\CloudDrive")) //预览目录
	err := r.Run(":5560")
	if err != nil {
		return
	}
}

// 客户端可以调用这个接口获取文件列表，包括文件夹和文件
func fileList() (string, func(c *gin.Context)) {
	return "fileList", func(c *gin.Context) {
		var fileList []FileBean
		dirs := []string{"D:\\temp"}
		for _, dir := range dirs {
			file, err := os.Open(dir)
			if err != nil {
				fmt.Printf("err：%s\n", err)
			}
			fileList = append(fileList, FileBean{FILENAME: dir, ISDIR: true, CHILDS: parseChildFile(*file)})
			//err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			//	if err != nil {
			//		return err
			//	}
			//
			//	// 判断是否是文件夹
			//	if info.IsDir() {
			//		fmt.Printf("目录：%s\n", path)
			//	} else {
			//		fmt.Printf("文件：%s\n", path)
			//	}
			//	fileBean := FileBean{FILENAME: info.Name(), ISDIR: info.IsDir()}
			//	fileList = append(fileList, fileBean)
			//
			//	return nil
			//})
			//
			//if err != nil {
			//	fmt.Println("遍历目录出错：", err)
			//}
		}
		//jsonData, err := json.Marshal(fileList)
		//if err != nil {
		//	fmt.Println("Error:", err)
		//	return
		//}

		responseSuccessWithContext(c, fileList)
	}
}

func parseChildFile(file os.File) []FileBean {
	result := make([]FileBean, 0)
	files, _ := file.ReadDir(-1)
	for _, child := range files {
		fileBean := FileBean{FILENAME: child.Name(), ISDIR: false}
		if child.IsDir() {
			fileBean.ISDIR = true
			childFile, _ := os.Open(file.Name() + "\\" + child.Name())
			fileBean.CHILDS = parseChildFile(*childFile)
		}
		result = append(result, fileBean)
	}
	return result
}

func apks() (string, func(c *gin.Context)) {
	return "apks", func(c *gin.Context) {
		http.StripPrefix("/files/", http.FileServer(http.Dir("D://temp"))).ServeHTTP(c.Writer, c.Request)
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
		androidDir, _ := os.Create("/data/data/com.termux/files/usr/share/nginx/html/upload/")
		_, err := androidDir.Stat()
		if err == nil {
			path = "/data/data/com.termux/files/usr/share/nginx/html/upload/"
		} else {
			androidDir, _ := os.Create("../usr/share/nginx/html/upload/")
			_, err := androidDir.Stat()
			if err == nil {
				path = "../usr/share/nginx/html/upload/"
			}
		}
		err = c.SaveUploadedFile(f, path+f.Filename)
		if err != nil {
			return
		}
		responseSuccessWithContext(c, "Upload File Success")
	}
}

func responseSuccess(data any) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": "0",
			"msg":  "OK",
			"data": data,
		})
	}
}
func responseSuccessWithContext(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
		"data": data,
	})
}

// 注意变量命名要大写，否则转json不能识别该变量
type FileBean struct {
	FILENAME string
	ISDIR    bool
	CHILDS   []FileBean
}
