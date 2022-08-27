package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/15110102/phuongpt3-market-server/src/app"
	"github.com/15110102/phuongpt3-market-server/src/model"
	"github.com/gorilla/mux"
)

func main() {
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
	var d app.AppIface = app.App{}
	result, err := d.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	vars := mux.Vars(r)
	productId := vars["id"]

	result, err := d.GetProduct(productId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	var product *model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := d.CreateProduct(product)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	var product *model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]
	product.Id = productId

	result, err := d.UpdateProduct(product)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	vars := mux.Vars(r)
	productId := vars["id"]

	result, err := d.DeleteProduct(productId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	var order *model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := d.CreateOrder(order)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	vars := mux.Vars(r)
	orderId := vars["id"]

	result, err := d.GetOrder(orderId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func updateOrderCallback(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	var cbOrder *model.CallbackOrder
	err := json.NewDecoder(r.Body).Decode(&cbOrder)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := d.UpdateOrderCallback(cbOrder)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getOrderStatusInThirdPartyServer(w http.ResponseWriter, r *http.Request) {
	var d app.AppIface = app.App{}
	vars := mux.Vars(r)
	appTransId := vars["appTransId"]

	result, err := d.GetOrderStatusInThirdPartyServer(appTransId)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}
