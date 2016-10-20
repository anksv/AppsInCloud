package routers

import (
	"../controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/marathon/CreateApp", &controllers.MarathonController{}, "get,post:CreateApp")
	//beego.Router("/marathon/DeleteApp", &controllers.MarathonController{}, "post:DeleteApp")
	//beego.Router("/marathon/ScaleApp", &controllers.MarathonController{}, "put:ScaleApp")
	//beego.Router("/marathon/ListApp", &controllers.MarathonController{}, "get:ListAPP")
	beego.Router("/v1/CreateWorkLoad", &controllers.CHEController{}, "post:CreateApp")
	beego.Router("/v1/DeleteWorkLoad", &controllers.CHEController{}, "post:DeleteApp")
	//beego.Router("/v1/ScaleWorkLoad", &controllers.CHEController{}, "put:ScaleApp")
	beego.Router("/v1/ListWorkLoad", &controllers.CHEController{}, "get:ListAPP")

}
