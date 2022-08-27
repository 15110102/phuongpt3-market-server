package model

type Product struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	CreateAt int64  `json:"created_at"`
	UpdateAt int64  `json:"updated_at"`
	Price    int64  `json:"price"`
	Image    string `json:"image"`
}
