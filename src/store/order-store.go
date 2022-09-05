package store

import (
	"fmt"

	"github.com/15110102/phuongpt3-market-server/src/model"
)

func (s Store) CreateOrder(order *model.Order) (*model.Order, error) {
	query := fmt.Sprintf("INSERT INTO Orders(AppUser, AppTransId, ZpTransToken, Item, CreateAt, TotalPrice, Status) VALUES ('%s', '%s','%s', '%s', %d, %d, '%s');", order.AppUser, order.AppTransId, "", order.Item, order.CreateAt, order.TotalPrice, order.Status)
	res, err := db.Exec(query)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	order.Id = id
	return order, nil
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

func (s Store) UpdateZpTransTokenOrderById(orderId int64, zpTransToken string) (bool, error) {
	updateQuery := fmt.Sprintf("Update Orders Set ZpTransToken = '%s' Where Id = %d", zpTransToken, orderId)
	updateZpTransToken, err := db.Query(updateQuery)
	if err != nil {
		return false, err
	}
	defer updateZpTransToken.Close()
	return true, nil
}

func (s Store) GetOrder(orderId int64) (*model.Order, error) {
	var order model.Order
	err := db.QueryRow("SELECT * FROM Orders WHERE Id = ?", orderId).Scan(&order.Id, &order.AppUser, &order.AppTransId, &order.ZpTransToken, &order.Item, &order.CreateAt, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s Store) SearchOrders(searchOrders *model.SearchOrders) (*model.SearchOrdersResponse, error) {
	queryCout := `SELECT COUNT(*) FROM Orders`
	query := `SELECT * FROM Orders`
	queryDate := ``
	queryStatus := ``
	queryLimit := ``
	if searchOrders.FromDate > 0 && searchOrders.ToDate > 0 {
		queryDate = fmt.Sprintf("CreateAt >= %d AND CreateAt <= %d", searchOrders.FromDate, searchOrders.ToDate)
	} else if searchOrders.FromDate > 0 {
		queryDate = fmt.Sprintf("CreateAt >= %d", searchOrders.FromDate)
	} else if searchOrders.ToDate > 0 {
		queryDate = fmt.Sprintf("CreateAt <= %d", searchOrders.ToDate)
	} else {
		queryDate = ""
	}

	if searchOrders.Status != "" {
		queryStatus = fmt.Sprintf("Status = '%s'", searchOrders.Status)
	}

	if searchOrders.Limit >= 0 && searchOrders.Offset >= 0 {
		queryLimit = fmt.Sprintf(" Limit %d Offset %d", searchOrders.Limit, searchOrders.Offset)
	}

	if queryDate != "" && queryStatus != "" {
		query = query + " WHERE " + queryDate + " AND " + queryStatus + queryLimit + ";"
		queryCout = queryCout + " WHERE " + queryDate + " AND " + queryStatus + ";"
	} else {
		if queryDate == "" && queryStatus != "" {
			query = query + " WHERE " + queryStatus + queryLimit + ";"
			queryCout = queryCout + " WHERE " + queryStatus + ";"
		} else if queryStatus == "" && queryDate != "" {
			query = query + " WHERE " + queryDate + queryLimit + ";"
			queryCout = queryCout + " WHERE " + queryDate + ";"
		} else {
			query = query + queryLimit + ";"
		}
	}

	var count int64
	err := db.QueryRow(queryCout).Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.Id, &order.AppUser, &order.AppTransId, &order.ZpTransToken, &order.Item, &order.CreateAt, &order.TotalPrice, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	result := model.SearchOrdersResponse{
		Total: count,
		Data:  orders,
	}
	return &result, nil
}
