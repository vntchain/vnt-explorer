package routers

import (
	"github.com/vntchain/vnt-explorer/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/v1/test", &controllers.TestController{})

	beego.Router("/v1/test/:id", &controllers.TestController{})
}
