package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/15110102/phuongpt3-market-server/src/app"
	"github.com/15110102/phuongpt3-market-server/src/model"
	"github.com/15110102/phuongpt3-market-server/src/store"
	"github.com/gorilla/mux"
)

var a app.AppIface = app.App{}

func main() {
	store.InitDbConn() //DB Connection

	// API routes
	r := mux.NewRouter()
	//Domain product
	r.HandleFunc("/products", getAllProducts).Methods("GET")
	r.HandleFunc("/product", createProduct).Methods("POST")
	r.HandleFunc("/product/{id}", getProduct).Methods("GET")
	r.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")

	//Domain order
	r.HandleFunc("/order", createOrder).Methods("POST")
	r.HandleFunc("/order/callback", updateOrderCallback).Methods("POST")                              //API will call by ThirdPartyServer
	r.HandleFunc("/order/app-trans-id/{appTransId}", getOrderStatusInThirdPartyServer).Methods("GET") //--> Need remove and replace by cronjob
	r.HandleFunc("/order/{id}", getOrder).Methods("GET")                                              //--> Using by timmer in Merchant client

	http.Handle("/", r)
	port := ":5000"
	fmt.Println("Server is running on port" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	result, err := a.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	result, err := a.GetProduct(productId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product *model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := a.CreateProduct(product)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	var product *model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]
	product.Id = productId

	result, err := a.UpdateProduct(product)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	result, err := a.DeleteProduct(productId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order *model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := a.CreateOrder(order)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId := vars["id"]
	intVar, _ := strconv.ParseInt(orderId, 0, 64)
	result, err := a.GetOrder(intVar)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func updateOrderCallback(w http.ResponseWriter, r *http.Request) {
	var cbOrder *model.CallbackOrder
	err := json.NewDecoder(r.Body).Decode(&cbOrder)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := a.UpdateOrderCallback(cbOrder)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getOrderStatusInThirdPartyServer(w http.ResponseWriter, r *http.Request) {
	var a app.AppIface = app.App{}
	vars := mux.Vars(r)
	appTransId := vars["appTransId"]

	result, err := a.GetOrderStatusInThirdPartyServer(appTransId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}
