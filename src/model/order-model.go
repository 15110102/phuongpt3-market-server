package model

type Order struct {
	Id         string `json:"id"`
	AppUser    string `json:"app_user"`
	AppTransId string `json:"app_trans_id"`
	Item       string `json:"item"`
	CreateAt   int64  `json:"created_at"`
	TotalPrice int64  `json:"price"`
	Status     string `json:"status"`
}

type OrderToThirdPartyRequest struct {
	AppId       int64  `json:"app_id"`
	AppUser     string `json:"app_user"`
	AppTransId  string `json:"app_trans_id"`
	AppTime     int64  `json:"app_time"`
	Amount      string `json:"amount"`
	Item        string `json:"item"`
	Description string `json:"description"`
	EmbedDate   string `json:"embed_data"`
	BankCode    string `json:"bank_code"`
	Mac         string `json:"mac"`
	CallBackUrl string `json:"callback_url"`
	DeviceInfo  string `json:"device_info"`
	SubAppId    string `json:"sub_app_id"`
	Title       string `json:"title"`
	Currency    string `json:"currency"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

type OrderInThirdPartyResponse struct {
	OrderId          string `json:"order_id"`
	ReturnCode       int64  `json:"return_code"`
	ReturnMessage    string `json:"return_message"`
	SubReturnCode    int64  `json:"sub_return_code"`
	SubReturnMessage string `json:"sub_return_message"`
	OrderUrl         string `json:"order_url"`
	ZpTransToken     string `json:"zp_trans_token"`
}

type CallbackOrder struct {
	Data string `json:"data"`
	Mac  string `json:"mac"`
	Type int64  `json:"type"`
}

type CallbackOrderResponse struct {
	ReturnCode    int64  `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}

type CheckOrderStatusInThirdPartyResponse struct {
	ReturnCode       int64  `json:"return_code"`
	ReturnMessage    string `json:"return_message"`
	SubReturnCode    int64  `json:"sub_return_code"`
	SubReturnMessage string `json:"sub_return_message"`
	IsProcessing     bool   `json:"is_processing"`
	Amount           int64  `json:"amount"`
	ZpTransId        int64  `json:"zp_trans_id"`
}
