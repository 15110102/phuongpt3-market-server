package app

import "github.com/15110102/phuongpt3-market-server/src/model"

const (
	APP_ID = "2553"
	KEY_1  = "PcY4iZIKFCIdgZvA6ueMcMHHUbRLYjPL"
	KEY_2  = "kLtgPl8HHhfvMuDHPwKfgfsY4Ydm9eIz"
)

const MYSQL_CONNECTION_STRING = "root:Aa123456@tcp(localhost:3306)/market_server_db"

var PRODUCTS = []model.Product{
	{Id: "1", Name: "Macbook", Desc: "macbook pro 13 inch", CreateAt: 1661498415000, UpdateAt: 1661498415000, Price: 36000000, Image: ""},
	{Id: "2", Name: "Iphone", Desc: "Iphone 13", CreateAt: 1661498529000, UpdateAt: 1661498529000, Price: 21000000, Image: ""},
	{Id: "3", Name: "Imac", Desc: "Imac pro 2022", CreateAt: 1661498614000, UpdateAt: 1661498614000, Price: 42000000, Image: ""},
	{Id: "4", Name: "Applewatch", Desc: "applewatch seri 6", CreateAt: 1661498655000, UpdateAt: 1661498655000, Price: 8000000, Image: ""},
	{Id: "5", Name: "Airpod", Desc: "airpod 3", CreateAt: 1661498697000, UpdateAt: 1661498697000, Price: 4500000, Image: ""},
	{Id: "6", Name: "Airtag", Desc: "air tag 1", CreateAt: 1661580461000, UpdateAt: 1661580461000, Price: 900000, Image: ""},
	{Id: "7", Name: "Macpro", Desc: "Macpro", CreateAt: 1661498697000, UpdateAt: 1661498697000, Price: 125000000, Image: ""},
	{Id: "8", Name: "MacMini", Desc: "Macmini", CreateAt: 1661580592000, UpdateAt: 1661580592000, Price: 28000000, Image: ""},
	{Id: "9", Name: "Macos", Desc: "Mac os", CreateAt: 1661580592000, UpdateAt: 1661580592000, Price: 3000000, Image: ""},
	{Id: "10", Name: "MacAir", Desc: "airpod 3", CreateAt: 1661498697000, UpdateAt: 1661498697000, Price: 25000000, Image: ""},
}

const (
	NEW     = "new"
	FAILED  = "failed"
	SUCCESS = "success"
)

const DOMAIN_THIRD_PARTY = "https://sb-openapi.zalopay.vn"
const DOMAIN_API = "https://0ee8-2401-d800-584a-6e01-a932-4ec-a885-6ca6.ap.ngrok.io"
