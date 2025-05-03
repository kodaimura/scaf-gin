package input

type Account struct {
	AccountId       int
	AccountName     string
	AccountPassword string
}

type AccountPK struct {
	AccountId int
}

type Login struct {
	AccountName     string
	AccountPassword string
}

type UpdatePassword struct {
	AccountId       int
	AccountPassword string
}
