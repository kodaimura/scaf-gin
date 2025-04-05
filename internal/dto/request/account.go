package request

type Signup struct {
	AccountName     string `json:"account_name"`
	AccountPassword string `json:"account_password"`
}

type Login struct {
	AccountName     string `json:"account_name"`
	AccountPassword string `json:"account_password"`
}

type PutAccount struct {
	AccountName string `json:"account_name"`
}

type PutAccountPassword struct {
	OldAccountPassword string `json:"old_account_password"`
	NewAccountPassword string `json:"new_account_password"`
}