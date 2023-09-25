package app

type UserFormLogin struct {
	Email    string `json:"email"  valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

type UserFormRegister struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

type UserFormUpdate struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password"`
}

type UserResult struct {
	ID       uint    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
