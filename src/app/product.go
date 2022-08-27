package app

import (
	"github.com/15110102/phuongpt3-market-server/src/model"
)

func (a App) GetAllProducts() ([]model.Product, error) {
	//TODO: Check Valid params
	var products = PRODUCTS

	return products, nil
}

func (a App) GetProduct(productId string) (*model.Product, error) {
	//TODO: Check Valid params
	var products = PRODUCTS
	for _, product := range products {
		if product.Id == productId {
			return &product, nil
		}
	}
	return nil, nil
}

func (a App) CreateProduct(product *model.Product) (*model.Product, error) {
	//TODO: Check Valid params
	return product, nil
}

func (a App) UpdateProduct(product *model.Product) (*model.Product, error) {
	//TODO: Check Valid params
	return product, nil
}

func (a App) DeleteProduct(product string) (bool, error) {
	//TODO: Check Valid params
	return true, nil
}
