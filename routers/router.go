package routers

import (
	"github.com/vntchain/vnt-explorer/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})

	beego.Router("/v1/blocks", &controllers.BlockController{}, "get:List")
	beego.Router("/v1/blocks/count", &controllers.BlockController{}, "get:Count")
	beego.Router("/v1/block/:n_or_h", &controllers.BlockController{})

	beego.Router("/v1/txs", &controllers.TransactionController{}, "get:List;post:Post")
}
