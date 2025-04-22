package input

type Account struct {
	AccountId       int
	AccountName     string
	AccountPassword string
}

type AccountPK struct {
	AccountId int
}

type Signup struct {
	AccountName     string
	AccountPassword string
}

type Login struct {
	AccountName     string
	AccountPassword string
}
