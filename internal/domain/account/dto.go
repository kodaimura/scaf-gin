package account

type GetDto struct {
	Id   int
	Name string
}

type GetOneDto struct {
	Id int
}

type UpdateOneDto struct {
	Id   int
	Name string
}

type DeleteOneDto struct {
	Id int
}
