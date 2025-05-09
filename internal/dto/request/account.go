package request

type Signup struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type PutAccount struct {
	Name string `json:"name" binding:"required"`
}

type PutPassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
