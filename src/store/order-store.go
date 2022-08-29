package store

import (
	"fmt"

	"github.com/15110102/phuongpt3-market-server/src/model"
)

func (s Store) CreateOrder(order *model.Order) (bool, error) {
	query := fmt.Sprintf("INSERT INTO Orders VALUES ( '%s', '%s', '%s', '%s', %d, %d, '%s')", order.Id, order.AppUser, order.AppTransId, order.Item, order.CreateAt, order.TotalPrice, order.Status)
	insert, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer insert.Close()

	return true, nil
}

func (s Store) UpdateStatusOrderByTrans(transId string, status string) (bool, error) {
	updateQuery := fmt.Sprintf("Update Orders Set Status = '%s' Where AppTransId = '%s'", status, transId)
	updateStatus, err := db.Query(updateQuery)
	if err != nil {
		return false, err
	}
	defer updateStatus.Close()
	return true, nil
}

func (s Store) GetOrder(orderId string) (*model.Order, error) {
	var order model.Order
	err := db.QueryRow("SELECT * FROM Orders WHERE Id = ?", orderId).Scan(&order.Id, &order.AppUser, &order.AppTransId, &order.Item, &order.CreateAt, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
