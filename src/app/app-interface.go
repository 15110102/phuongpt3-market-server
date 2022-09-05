package app

import "github.com/15110102/phuongpt3-market-server/src/model"

type App struct{}

type AppIface interface {
	GetAllProducts() ([]model.Product, error)
	GetProduct(productId string) (*model.Product, error)
	CreateProduct(product *model.Product) (*model.Product, error)
	UpdateProduct(product *model.Product) (*model.Product, error)
	DeleteProduct(productId string) (bool, error)
	CreateOrder(order *model.Order) (*model.OrderInThirdPartyResponse, error)
	GetOrder(orderId int64) (*model.Order, error)
	UpdateOrderCallback(cbOrder *model.CallbackOrder) (*model.CallbackOrderResponse, error)
	GetOrderStatusInThirdPartyServer(appTransId string) (*model.CheckOrderStatusInThirdPartyResponse, error)
}
