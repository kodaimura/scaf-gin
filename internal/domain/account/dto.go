package account

type GetDto struct {
	Id   int
	Name string
}

type GetOneDto struct {
	Id   int
}

type CreateOneDto struct {
	Name     string
	Password string
}

type UpdateOneDto struct {
	Id int
	Name     string
}

type DeleteOneDto struct {
	Id int
}

type LoginDto struct {
	Name     string
	Password string
}

type UpdatePasswordDto struct {
	Id       int
	Password string
}
