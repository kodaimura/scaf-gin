package account

import accountModule "scaf-gin/internal/module/account"

type ListDto struct{}

func (uc *usecase) List(in ListDto) ([]accountModule.Account, error) {
	return uc.accountModule.GetAll()
}
