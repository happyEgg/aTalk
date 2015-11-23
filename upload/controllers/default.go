package controllers

import (
	"fmt"
	"io/ioutil"
	"os"

	"time"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// @router / [get]
func (c *MainController) LoadUpload() {
	c.TplNames = "upload.html"
}

// @router /download/* [get]
func (c *MainController) Download() {
	filename := c.Ctx.Input.Param(":splat")
	beego.Debug(filename)
	dstfilename := "./static/" + filename
	c.Ctx.Output.Download(dstfilename, filename)
}

// @router /uploadform [post]
func (c *MainController) UploadForm() {

	_, file, err := c.GetFile("filename")
	if nil == err {
		if serr := c.SaveToFile("filename", "./static/"+file.Filename); serr == nil {
		} else {
			beego.Error(serr)
			c.Ctx.WriteString(serr.Error())
		}
		c.Ctx.ResponseWriter.Write([]byte("http://localhost:8989/download/" + file.Filename))
		return
	}
	beego.Error(err)
	c.Ctx.WriteString(err.Error())
}

// @router /upload/* [*]
func (c *MainController) Upload() {

	rw := c.Ctx.ResponseWriter
	req := c.Ctx.Request
	if req.Method == "GET" {
		rw.Write([]byte(""))
	}
	req.ParseForm()
	length := req.Header.Get("Content-Length")
	fmt.Println(length)
	b, err := ioutil.ReadAll(req.Body)
	if checkerr(err) {
		rw.Write([]byte("error"))
	}
	filename := c.Ctx.Input.Param(":splat")
	beego.Debug(filename)
	file, err := os.OpenFile("./static/"+filename, os.O_CREATE|os.O_WRONLY, 0644)
	if checkerr(err) {
		rw.Write([]byte("error"))
	}
	_, err = file.Write(b)
	if checkerr(err) {
		rw.Write([]byte("error"))
	}
	rw.Write([]byte("http://localhost:8989/download/" + filename))
}

func checkerr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func randomFilename() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
