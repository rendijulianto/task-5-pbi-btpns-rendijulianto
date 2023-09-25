package app

type PhotoFormCreate struct {
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required"`
}

type PhotoFormUpdate struct {
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required"`
}

type PhotoResult struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email   string `json:"email"`
}
