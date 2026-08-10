package account

import accountModule "scaf-gin/internal/module/account"

type Usecase interface {
	List(in ListDto) ([]accountModule.Account, error)
	Get(in GetDto) (accountModule.Account, error)
	Create(in CreateDto) (accountModule.Account, error)
	Update(in UpdateDto) (accountModule.Account, error)
	Disable(in DisableDto) (accountModule.Account, error)
	Enable(in EnableDto) (accountModule.Account, error)
	Delete(in DeleteDto) error
}

type usecase struct {
	accountModule accountModule.Module
}

func NewUsecase(accountModule accountModule.Module) Usecase {
	return &usecase{
		accountModule: accountModule,
	}
}
