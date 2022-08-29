package store

import (
	"database/sql"
	"fmt"

	"github.com/15110102/phuongpt3-market-server/src/model"
	_ "github.com/go-sql-driver/mysql"
)

type Store struct{}
type StoreIface interface {
	CreateOrder(order *model.Order) (bool, error)
	GetOrder(orderId string) (*model.Order, error)
	UpdateStatusOrderByTrans(transId string, status string) (bool, error)
}

var db *sql.DB

func InitDbConn() {
	var err error
	db, err = sql.Open("mysql", MYSQL_CONNECTION_STRING)
	if err != nil {
		fmt.Println(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(err)
	}
}
