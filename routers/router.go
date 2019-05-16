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
	beego.Router("/v1/txs/history", &controllers.TransactionController{}, "get:History")
	beego.Router("/v1/tx/:tx_hash", &controllers.TransactionController{})

	beego.Router("/v1/accounts", &controllers.AccountController{}, "get:List;post:Post")
	beego.Router("/v1/accounts/count", &controllers.AccountController{}, "get:Count")
	beego.Router("/v1/account/:address", &controllers.AccountController{})
	beego.Router("/v1/account/:address/tokens", &controllers.TokenBalanceController{}, "get:ListByAccount")
	beego.Router("/v1/account/:address/tokens/count", &controllers.TokenBalanceController{}, "get:TokenCount")

	beego.Router("/v1/nodes", &controllers.NodeController{}, "get:List;post:Post")
	beego.Router("/v1/nodes/count", &controllers.NodeController{}, "get:Count")
	beego.Router("/v1/node/:address", &controllers.NodeController{})

	beego.Router("/v1/token/:address/holders", &controllers.TokenBalanceController{}, "get:ListByToken")
	beego.Router("/v1/token/:address/holders/count", &controllers.TokenBalanceController{}, "get:HolderCount")

	beego.Router("/v1/stats", &controllers.NetController{}, "get:Stats")

	beego.Router("/v1/search/:keyword", &controllers.SearchController{}, "get:Search")

	beego.Router("/v1/hydrant", &controllers.HydrantController{}, "post:SendVnt")

	beego.Router("/v1/kline",&controllers.MarketController{},"get:History")
	beego.Router("/v1/market",&controllers.MarketController{},"get:Market")
}
