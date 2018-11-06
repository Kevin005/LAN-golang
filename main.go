package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"net"
)

var htmlLAN = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>局域网内文件互传</title>
</head>
<body style="text-align: center;"> 
    <h1>局域网内文件互传</h1>
    <br>
    <br>
    <form action="UploadFile" method="post" enctype="multipart/form-data">
    <input type="file" name="fileUpload" />
    <input type="submit" name="上传文件">
    </form>
        <br>
    <br>
        <br>
    <br>
    <a href="/file">文件下载</a>
</body>
</html>
`

func main() {
	fmt.Println("please visit the link below")
	showip()
	//可以下载到目录：/Users/kevin/testLAN/
	http.HandleFunc("/", uploadFileHandler)
	//todo 只是列出了文件
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir("/Users/kevin/temp"))))
	http.ListenAndServe(":8080", nil)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, htmlLAN)
	if r.Method == "POST" {
		file, handler, err := r.FormFile("fileUpload") //name的字段
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		check(err)
		newFile, err := os.Create("/Users/kevin/testLAN/" + handler.Filename)
		check(err)
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			check(err)
			return
		}
		fmt.Println(" upload successfully:" + "/Users/kevin/testLAN/" + handler.Filename)
		w.Write([]byte("<br/> <br/>download success " + handler.Filename + "<br/>please check path /Users/kevin/testLAN/"))
	}
}

func check(e error) bool {
	if e != nil {
		fmt.Println(e)
		return true
	}
	return false
}

func showip() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String() + ":8080")
			}
		}
	}
}
