package common

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	ERROR_SYSTEM 			 = "system_err"			//系统错误
	ERROR_WRONG_ADDRESS 	 = "wrong_address"		//地址格式错误
	ERROR_DUPLICATED_SEND 	 = "duplicated_send"	//正在发送中，请稍后查看
	ERROR_SEND_TO_FREQUENTLY = "send_over_frequent"	//发送太频繁，一个小时只能发送一次
	ERROR_NONCE_ERROR 		 = "system_nonce_err"	//获取交易nonce失败，请联系管理员
	ERROR_SIGN_ERROR 		 = "system_sign_err"	//交易签名失败，请联系管理员
	ERROR_TX_SEND_ERROR 	 = "system_tx_send_err"	//交易发送失败，请联系管理员
	ERROR_WRONG_KEYWORD 	 = "wrong_keyword"		//关键字格式错误：请输入合法的区块hash，区块号，交易hash或账户地址进行搜索
	ERROR_SEARCH_ERROR 		 = "search_err"			//搜索错误
	ERROR_NOT_FOUND 		 = "not_found"			//不存在
)