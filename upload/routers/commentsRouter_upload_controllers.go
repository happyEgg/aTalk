package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["upload/controllers:MainController"] = append(beego.GlobalControllerRouter["upload/controllers:MainController"],
		beego.ControllerComments{
			"LoadUpload",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["upload/controllers:MainController"] = append(beego.GlobalControllerRouter["upload/controllers:MainController"],
		beego.ControllerComments{
			"Download",
			`/download/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["upload/controllers:MainController"] = append(beego.GlobalControllerRouter["upload/controllers:MainController"],
		beego.ControllerComments{
			"UploadForm",
			`/uploadform`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["upload/controllers:MainController"] = append(beego.GlobalControllerRouter["upload/controllers:MainController"],
		beego.ControllerComments{
			"Upload",
			`/upload/*`,
			[]string{"*"},
			nil})

}
