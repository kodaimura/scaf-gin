package response

type Account struct {
	AccountId   int    `json:"account_id"`
	AccountName string `json:"account_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Login struct {
	AccessToken      string  `json:"access_token"`
	RefreshToken     string  `json:"refresh_token"`
	AccessExpiresIn  int     `json:"access_expires_in"`
	RefreshExpiresIn int     `json:"refresh_expires_in"`
	Account          Account `json:"account"`
}

type Refresh struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
