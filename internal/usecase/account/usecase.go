package account

import (
	accountModule "scaf-gin/internal/module/account"

	"gorm.io/gorm"
)

type Usecase interface {
	Get(in GetDto) ([]accountModule.Account, error)
	GetOne(in GetOneDto) (accountModule.Account, error)
	UpdateOne(in UpdateOneDto) (accountModule.Account, error)
	DeleteOne(in DeleteOneDto) error
}

type usecase struct {
	db             *gorm.DB
	accountService accountModule.Service
}

func NewUsecase(
	db *gorm.DB,
	accountService accountModule.Service,
) Usecase {
	return &usecase{
		db:             db,
		accountService: accountService,
	}
}

func (uc *usecase) Get(in GetDto) ([]accountModule.Account, error) {
	return uc.accountService.Get(accountModule.Account{}, uc.db)
}

func (uc *usecase) GetOne(in GetOneDto) (accountModule.Account, error) {
	return uc.accountService.GetOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
}

func (uc *usecase) UpdateOne(in UpdateOneDto) (accountModule.Account, error) {
	return uc.accountService.UpdateOne(accountModule.Account{
		Id:   in.Id,
		Name: in.Name,
	}, uc.db)
}

func (uc *usecase) DeleteOne(in DeleteOneDto) error {
	return uc.accountService.DeleteOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
}
