package request

type Signup struct {
	AccountName     string `json:"account_name" binding:"required"`
	AccountPassword string `json:"account_password" binding:"required,min=8"`
}

type Login struct {
	AccountName     string `json:"account_name"`
	AccountPassword string `json:"account_password"`
}

type PutAccount struct {
	AccountName string `json:"account_name" binding:"required"`
}

type PutAccountPassword struct {
	OldAccountPassword string `json:"old_account_password" binding:"required"`
	NewAccountPassword string `json:"new_account_password" binding:"required"`
}
