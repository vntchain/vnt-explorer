package models

type Account struct {
	Id             int
	Address        string
	Vname          string
	Balance        string
	TxCount        int
	IsContract     bool
	ContractName   string
	ContractOwner  string
	Code           string
	Abi            string
	Home           string
	InitTx         string
	IsToken        bool
	TokenType      int
	TokenSymbol    string
	TokenLogo      string
	TokenAmount    string
	TokenAcctCount string
	FirstBlock     string
	LastBlock      string
}
