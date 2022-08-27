package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/15110102/phuongpt3-market-server/src/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/zpmep/hmacutil"
)

func (a App) CreateOrder(order *model.Order) (*model.OrderToThirdPartyResponse, error) {
	//TODO: Check Valid input params
	type object map[string]interface{}
	var (
		app_id = APP_ID
		key1   = KEY_1
	)
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000)
	embedData, _ := json.Marshal(object{})
	params := make(url.Values)
	params.Add("app_id", app_id)
	params.Add("amount", fmt.Sprintf("%d", order.TotalPrice))
	params.Add("app_user", order.AppUser)
	params.Add("embed_data", string(embedData))
	params.Add("item", order.Item)
	params.Add("description", "PhuongPT - Payment test for the order #"+strconv.Itoa(transID))
	params.Add("bank_code", "zalopayapp")

	now := time.Now()
	timeNowMil := now.UnixNano() / int64(time.Millisecond)
	params.Add("app_time", strconv.FormatInt(timeNowMil, 10))

	appTransId := fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID)
	fmt.Println("appTransId: ", appTransId)
	params.Add("app_trans_id", appTransId)
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var orderToThirdPartyResponse *model.OrderToThirdPartyResponse
	err = json.Unmarshal(responseData, &orderToThirdPartyResponse)
	if err != nil {
		return nil, err
	}

	//TODO:Cronjob check status order: GetOrderStatusInThirdPartyServer
	db, err := sql.Open("mysql", MYSQL_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	id := uuid.New()
	query := fmt.Sprintf("INSERT INTO Orders VALUES ( '%s', '%s', '%s', '%s', %d, %d, '%s')", id.String(), order.AppUser, appTransId, order.Item, timeNowMil, order.TotalPrice, "New")
	orderToThirdPartyResponse.OrderId = id.String()
	insert, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return orderToThirdPartyResponse, nil
}

func (a App) GetOrder(orderId string) (*model.Order, error) {
	//TODO: Check Valid input params
	db, err := sql.Open("mysql", MYSQL_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var order model.Order
	err = db.QueryRow("SELECT * FROM Orders WHERE Id = ?", orderId).Scan(&order.Id, &order.AppUser, &order.AppTransId, &order.Item, &order.CreateAt, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (a App) UpdateOrderCallback(cbOrder *model.CallbackOrder) (*model.CallbackOrderResponse, error) {
	//TODO: Check Valid input params
	var key2 = KEY_2
	requestMac := cbOrder.Mac
	dataStr := cbOrder.Data
	mac := hmacutil.HexStringEncode(hmacutil.SHA256, key2, dataStr)

	var result model.CallbackOrderResponse
	// mac := "d8d33baf449b31d7f9b94fa50d7c942c08cd4d83f28fa185557da21acb104f67"
	if mac != requestMac {
		result.ReturnCode = -1
		result.ReturnMessage = "mac not equal"
		return nil, nil
	} else {
		result.ReturnCode = 1
		result.ReturnMessage = "success"

		// TODO: merchant cập nhật trạng thái cho đơn hàng
		var dataJSON map[string]interface{}
		json.Unmarshal([]byte(dataStr), &dataJSON)
		fmt.Println("update order's status = success where app_trans_id =", dataJSON["app_trans_id"])

		return &result, nil
	}
}

func (a App) GetOrderStatusInThirdPartyServer(appTransId string) (*model.CheckOrderStatusInThirdPartyResponse, error) {
	//TODO: Check Valid input params
	var (
		appID = APP_ID
		key1  = KEY_1
		// key2       = KEY_2
		appTransID = appTransId
	)

	data := fmt.Sprintf("%v|%s|%s", appID, appTransID, key1)
	params := map[string]interface{}{
		"app_id":       appID,
		"app_trans_id": appTransID,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, key1, data),
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	responseData, _ := ioutil.ReadAll(res.Body)

	var orderStatus *model.CheckOrderStatusInThirdPartyResponse
	err = json.Unmarshal(responseData, &orderStatus)
	if err != nil {
		return nil, err
	}

	return orderStatus, err
}
