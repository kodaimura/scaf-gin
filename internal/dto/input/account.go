package input

type Account struct {
	Id       int
	Name     string
	Password string
}

type AccountPK struct {
	Id int
}

type Login struct {
	Name     string
	Password string
}

type UpdatePassword struct {
	Id       int
	Password string
}
