package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Order struct {
	ID           string  `json:"id"`
    Status       string  `json:"status"`
    Items        []Item  `json:"items"`
    Total        float64 `json:"total"`
    CurrencyUnit string  `json:"currencyUnit"`
}

type Item struct {
    ID          string  `json:"id"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Quantity    int     `json:"quantity"`
}

var orders []Order


func main() {
	router := mux.NewRouter()


	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", updateOrderStatus).Methods("PUT")
	router.HandleFunc("/orders", getOrders).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createOrder(w http.ResponseWriter, r *http.Request){
	var order Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	order.ID = uuid.NewString()
	orders = append(orders, order)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]


	var orderStatus struct {
		Status string `json:"status"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&orderStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    for i, order := range orders {
        if order.ID == id {
            orders[i].Status = orderStatus.Status
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(orders[i])
            return
        }
    }
    http.Error(w, fmt.Sprintf("order with ID %s not found", id), http.StatusNotFound)
}


// getOrders returns a list of orders sorted and filtered by various fields
func getOrders(w http.ResponseWriter, r *http.Request) {
  	// Parse query parameters
	query := r.URL.Query()

	// Filter orders by ID
	if ids, ok := query["id"]; ok {
		filteredOrders := make([]Order, 0)
		for _, order := range orders {
			for _, id := range ids {
				if order.ID == id {
					filteredOrders = append(filteredOrders, order)
					break
				}
			}
		}
		orders = filteredOrders
	}

	// Filter orders by status
	if statuses, ok := query["status"]; ok {
		filteredOrders := make([]Order, 0)
		for _, order := range orders {
			for _, status := range statuses {
				if order.Status == status {
					filteredOrders = append(filteredOrders, order)
					break
				}
			}
		}
		orders = filteredOrders
	}


	sortBy := query.Get("sortBy")
	switch sortBy {
	case "id":
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].ID < orders[j].ID
		})
	case "status":
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Status < orders[j].Status
		})
	case "total":
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].Total < orders[j].Total
		})
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}