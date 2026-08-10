package account

import (
	"scaf-gin/internal/model"
	"scaf-gin/internal/module"
)

type Usecase interface {
	List(in ListDto) ([]model.Account, error)
	Get(in GetDto) (model.Account, error)
	Create(in CreateDto) (model.Account, error)
	Update(in UpdateDto) (model.Account, error)
	Disable(in DisableDto) (model.Account, error)
	Enable(in EnableDto) (model.Account, error)
	Delete(in DeleteDto) error
}

type usecase struct {
	accountModule module.AccountModule
}

func NewUsecase(accountModule module.AccountModule) Usecase {
	return &usecase{
		accountModule: accountModule,
	}
}
