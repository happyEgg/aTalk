package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"],
		beego.ControllerComments{
			"LoadUpload",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"],
		beego.ControllerComments{
			"Download",
			`/download/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"],
		beego.ControllerComments{
			"UploadForm",
			`/uploadform`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/everfore/upload/controllers:MainController"],
		beego.ControllerComments{
			"Upload",
			`/upload/*`,
			[]string{"*"},
			nil})

}
