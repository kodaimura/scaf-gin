package account

type PutMeRequest struct {
	Name string `json:"name" binding:"required"`
}
