package app

import (
	"bytes"
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
	"github.com/15110102/phuongpt3-market-server/src/store"
	"github.com/google/uuid"
	"github.com/zpmep/hmacutil"
)

var s store.StoreIface = store.Store{}

func (a App) CreateOrder(order *model.Order) (*model.OrderInThirdPartyResponse, error) {
	//TODO: Check Valid input params
	now := time.Now()
	timeNowMil := now.UnixNano() / int64(time.Millisecond)
	order.CreateAt = timeNowMil
	transId := uuid.New()
	appTransId := fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transId.String()[:len(transId.String())-5])
	order.AppTransId = appTransId
	order.Status = NEW

	resultDb, err := s.CreateOrder(order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	orderToThirdPartyResponse, err := a.createOrderInThirdParty(order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	orderToThirdPartyResponse.OrderId = resultDb.Id

	_, err = s.UpdateZpTransTokenOrderById(orderToThirdPartyResponse.OrderId, orderToThirdPartyResponse.ZpTransToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Loop check status order every 5s and end check after 15m or give callback from third party server
	timeLoop := time.NewTicker(5 * time.Second)
	endTrackOrder := make(chan bool)
	go func() {
		time.Sleep(15 * time.Minute)
		endTrackOrder <- true
	}()

	go func(appTransId string, orderId int64) {
		for {
			select {
			//End after 15m
			case <-endTrackOrder:
				order, err := a.GetOrder(orderId)
				if err != nil {
					fmt.Println(err)
					return
				}
				if order.Status == SUCCESS {
					return
				} else {
					_, err := s.UpdateStatusOrderByTrans(appTransId, FAILED)
					if err != nil {
						fmt.Println(err)
						return
					}
					return
				}
			//Loop every 5s
			case <-timeLoop.C:
				order, err := a.GetOrder(orderId)
				if err != nil {
					fmt.Println(err)
					return
				}
				if order.Status == SUCCESS {
					return
				} else {
					result, err := a.GetOrderStatusInThirdPartyServer(appTransId)
					if err != nil {
						fmt.Println(err)
						return
					}
					//reupdate status order if third party server missing callback
					if result.ReturnCode == 1 {
						_, err := s.UpdateStatusOrderByTrans(appTransId, SUCCESS)
						if err != nil {
							fmt.Println(err)
							return
						}
						return
					}
				}
			}
		}
	}(appTransId, orderToThirdPartyResponse.OrderId)

	return orderToThirdPartyResponse, nil
}

func (a App) createOrderInThirdParty(order *model.Order) (*model.OrderInThirdPartyResponse, error) {
	type object map[string]interface{}
	var (
		app_id = APP_ID
		key1   = KEY_1
	)

	rand.Seed(time.Now().UnixNano())
	calbackUrl := fmt.Sprintf("%s/order/callback", DOMAIN_API)
	embedData, _ := json.Marshal(object{"redirecturl": DOMAIN_CLIENT})
	params := make(url.Values)
	params.Add("app_id", app_id)
	params.Add("amount", fmt.Sprintf("%d", order.TotalPrice))
	params.Add("app_user", order.AppUser)
	params.Add("callback_url", calbackUrl)
	params.Add("embed_data", string(embedData))
	params.Add("item", order.Item)
	params.Add("description", "PhuongPT - Payment test for the order ")
	params.Add("bank_code", "zalopayapp")
	params.Add("app_time", strconv.FormatInt(order.CreateAt, 10))
	params.Add("app_trans_id", order.AppTransId)
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	domain := fmt.Sprintf("%s/v2/create", DOMAIN_THIRD_PARTY)
	res, err := http.PostForm(domain, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var orderInThirdPartyResponse *model.OrderInThirdPartyResponse
	err = json.Unmarshal(responseData, &orderInThirdPartyResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return orderInThirdPartyResponse, nil
}

func (a App) GetOrder(orderId int64) (*model.Order, error) {
	//TODO: Check Valid input params
	order, err := s.GetOrder(orderId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return order, nil
}

func (a App) UpdateOrderCallback(cbOrder *model.CallbackOrder) (*model.CallbackOrderResponse, error) {
	//TODO: Check Valid input params
	var key2 = KEY_2
	requestMac := cbOrder.Mac
	dataStr := cbOrder.Data
	mac := hmacutil.HexStringEncode(hmacutil.SHA256, key2, dataStr)

	var result model.CallbackOrderResponse
	if mac != requestMac {
		result.ReturnCode = -1
		result.ReturnMessage = "mac not equal"
		return nil, nil
	} else {
		result.ReturnCode = 1
		result.ReturnMessage = "success"

		var dataJSON map[string]interface{}
		json.Unmarshal([]byte(dataStr), &dataJSON)

		transId := fmt.Sprintf("%s", dataJSON["app_trans_id"])
		_, err := s.UpdateStatusOrderByTrans(transId, SUCCESS)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &result, nil
	}
}

func (a App) GetOrderStatusInThirdPartyServer(appTransId string) (*model.CheckOrderStatusInThirdPartyResponse, error) {
	//TODO: Check Valid input params
	convertAppId, _ := strconv.Atoi(APP_ID)
	var (
		appID = convertAppId
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
		fmt.Println(err)
		return nil, err
	}

	domain := fmt.Sprintf("%s/v2/query", DOMAIN_THIRD_PARTY)
	res, err := http.Post(domain, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	responseData, _ := ioutil.ReadAll(res.Body)

	var orderStatus *model.CheckOrderStatusInThirdPartyResponse
	err = json.Unmarshal(responseData, &orderStatus)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orderStatus, err
}

func (a App) SearchOrders(searchOrders *model.SearchOrders) (*model.SearchOrdersResponse, error) {
	orders, err := s.SearchOrders(searchOrders)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orders, nil
}
