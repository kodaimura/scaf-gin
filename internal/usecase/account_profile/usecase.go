package account_profile

import (
	profileModule "scaf-gin/internal/module/account_profile"

	"gorm.io/gorm"
)

type Usecase interface {
	Get(in GetDto) ([]profileModule.AccountProfile, error)
	GetOne(in GetOneDto) (profileModule.AccountProfile, error)
	CreateOne(in CreateOneDto) (profileModule.AccountProfile, error)
	UpdateOne(in UpdateOneDto) (profileModule.AccountProfile, error)
	DeleteOne(in DeleteOneDto) error
}

type usecase struct {
	db                    *gorm.DB
	accountProfileService profileModule.Service
}

func NewUsecase(
	db *gorm.DB,
	accountProfileService profileModule.Service,
) Usecase {
	return &usecase{
		db:                    db,
		accountProfileService: accountProfileService,
	}
}

func (uc *usecase) Get(in GetDto) ([]profileModule.AccountProfile, error) {
	return uc.accountProfileService.Get(profileModule.AccountProfile{}, uc.db)
}

func (uc *usecase) GetOne(in GetOneDto) (profileModule.AccountProfile, error) {
	return uc.accountProfileService.GetOne(profileModule.AccountProfile{
		AccountId: in.AccountId,
	}, uc.db)
}

func (uc *usecase) CreateOne(in CreateOneDto) (profileModule.AccountProfile, error) {
	return uc.accountProfileService.CreateOne(profileModule.AccountProfile{
		AccountId:   in.AccountId,
		DisplayName: in.DisplayName,
		Bio:         in.Bio,
		AvatarURL:   in.AvatarURL,
	}, uc.db)
}

func (uc *usecase) UpdateOne(in UpdateOneDto) (profileModule.AccountProfile, error) {
	return uc.accountProfileService.UpdateOne(profileModule.AccountProfile{
		AccountId:   in.AccountId,
		DisplayName: in.DisplayName,
		Bio:         in.Bio,
		AvatarURL:   in.AvatarURL,
	}, uc.db)
}

func (uc *usecase) DeleteOne(in DeleteOneDto) error {
	return uc.accountProfileService.DeleteOne(profileModule.AccountProfile{
		AccountId: in.AccountId,
	}, uc.db)
}
