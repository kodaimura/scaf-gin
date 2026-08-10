package account

type GetDto struct {
	Id int
}

type GetOneDto struct {
	Id int
}

type CreateOneDto struct {
	LoginID   *string
	Email     *string
	Password  string
	FirstName string
	LastName  string
}

type UpdateOneDto struct {
	Id        int
	LoginID   *string
	Email     *string
	Password  *string
	FirstName string
	LastName  string
}

type DeleteOneDto struct {
	Id int
}

type DisableOneDto struct {
	Id int
}

type EnableOneDto struct {
	Id int
}
