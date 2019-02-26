package models

type Account struct {
	Address        string `orm:"pk"`
	Vname          string `orm:"unique"`
	Balance        string
	TxCount        uint64
	IsContract     bool
	ContractName   string
	ContractOwner  *Account `orm:"rel(fk)"`
	Code           string
	Abi            string
	Home           string
	InitTx         *Transaction `orm:"rel(fk)"`
	IsToken        bool
	TokenType      int
	TokenSymbol    string
	TokenLogo      string
	TokenAmount    string
	TokenAcctCount string
	FirstBlock     string
	LastBlock      string
}
