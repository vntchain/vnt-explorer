package routers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/controllers"
)

func init() {
	//beego.Router("/", &controllers.MainController{})

	beego.Router("/v1/blocks", &controllers.BlockController{}, "get:List;post:Post")
	beego.Router("/v1/blocks/count", &controllers.BlockController{}, "get:Count")
	beego.Router("/v1/block/:n_or_h", &controllers.BlockController{})

	beego.Router("/v1/txs", &controllers.TransactionController{}, "get:List;post:Post")
	beego.Router("/v1/txs/count", &controllers.TransactionController{}, "get:Count")
	beego.Router("/v1/tx/:tx_hash", &controllers.TransactionController{})

	beego.Router("/v1/accounts", &controllers.AccountController{}, "get:List;post:Post")
	beego.Router("/v1/accounts/count", &controllers.AccountController{}, "get:Count")
	beego.Router("/v1/account/:address", &controllers.AccountController{})
	beego.Router("/v1/account/:address/tokens", &controllers.TokenBalanceController{}, "get:ListByAccount")

	beego.Router("/v1/nodes", &controllers.NodeController{}, "get:List;post:Post")
	beego.Router("/v1/node/:address", &controllers.NodeController{})

	beego.Router("/v1/tokens/count", &controllers.TokenBalanceController{}, "get:CountByToken")
	beego.Router("/v1/token/:address/holders", &controllers.TokenBalanceController{}, "get:ListByToken")
}
