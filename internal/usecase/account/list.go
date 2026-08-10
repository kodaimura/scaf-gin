package account

import "scaf-gin/internal/model"

type ListDto struct{}

func (uc *usecase) List(in ListDto) ([]model.Account, error) {
	return uc.accountModule.GetAll()
}
