package routers

import (
	"github.com/vntchain/vnt-explorer/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/v1/test", &controllers.TestController{})

	beego.Router("/v1/test/:id", &controllers.TestController{})

	beego.Router("/v1/blocks", &controllers.BlockController{}, "get:List")

	beego.Router("/v1/block/:n_or_h", &controllers.BlockController{})
}
