package account_profile

type GetDto struct{}

type GetOneDto struct {
	AccountId int
}

type CreateOneDto struct {
	AccountId   int
	DisplayName string
	Bio         string
	AvatarURL   string
}

type UpdateOneDto struct {
	AccountId   int
	DisplayName string
	Bio         string
	AvatarURL   string
}

type DeleteOneDto struct {
	AccountId int
}
