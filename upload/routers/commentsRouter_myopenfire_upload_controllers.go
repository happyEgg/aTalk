package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"],
		beego.ControllerComments{
			"LoadUpload",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"],
		beego.ControllerComments{
			"Download",
			`/download/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"],
		beego.ControllerComments{
			"UploadForm",
			`/uploadform`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["myopenfire/upload/controllers:MainController"],
		beego.ControllerComments{
			"Upload",
			`/upload/*`,
			[]string{"*"},
			nil})

}
