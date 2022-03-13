package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/srrmendez/private-api-order/model"

	pkgHttp "github.com/srrmendez/services-interface-tools/pkg/http"
)

// healthCheck godoc
// @Tags HealthCheck
// @Accept  json
// @Produce json
// @Success 200
// @Router /health-check/ [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.JsonResponse(w, map[string]string{"status": "Running"}, http.StatusOK)
}

// Search Orders godoc
// @Tags Search Orders
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param status query string false "subscriber status"
// @Param category query string false "category"
// @Param order_type query string false "order type"
// @Success 200 {array} model.Order
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/ [get]
func searchOrders(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var status *model.OrderStatusType
	var category *model.CategoryType
	var orderType *model.OrderType

	if r.URL.Query().Get("status") != "" {
		st := model.OrderStatusType(r.URL.Query().Get("status"))
		status = &st
	}

	if r.URL.Query().Get("category") != "" {
		st := model.CategoryType(r.URL.Query().Get("category"))
		category = &st
	}

	if r.URL.Query().Get("order_type") != "" {
		st := model.OrderType(r.URL.Query().Get("order_type"))
		orderType = &st
	}

	orders, err := env.Services.orderService.Search(r.Context(), clientID, status, category, orderType)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, orders, http.StatusOK)
}

// Create Order godoc
// @Tags Create Order
// @Summary Create Order from Web Portal
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param req body model.Order true "Order"
// @Success 201 {object} model.Order
// @Failure 400 Incorrect body format
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/ [post]
func createOrder(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var order model.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = env.Validators.orderRequestValidator.Validate(order); err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	ord, err := env.Services.orderService.Create(r.Context(), order, clientID)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, ord, http.StatusCreated)
}

// Get Order godoc
// @Tags Get Order
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param id path string true "id"
// @Success 200 {object} model.Order
// @Failure 404 Order Not Found
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/{id} [get]
func getOrder(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	ord, err := env.Services.orderService.Get(r.Context(), id)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if ord == nil {
		pkgHttp.ErrorResponse(w, errors.New("order not found"), http.StatusNotFound)
		return
	}

	pkgHttp.JsonResponse(w, nil, http.StatusOK)
}

// Create Service Order godoc
// @Tags Create Service Order
// @Summary Create Service Order from commercial system
// @Accept  json
// @Produce json
// @Param x-client-id header string true "client id"
// @Param transaction_id query string true "transaction id"
// @Param req body model.ServiceOrderRequest true "Order"
// @Success 200 {object} model.ServiceOrderResponse
// @Failure 400 Incorrect body format
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/service [post]
func createServiceOrder(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	transactionID := r.URL.Query().Get("transaction_id")

	if transactionID == "" {
		pkgHttp.ErrorResponse(w, errors.New("external system transaction must be provided"), http.StatusBadRequest)
		return
	}

	// TODO Remove this
	var data interface{}
	json.NewDecoder(r.Body).Decode(&data)

	d, _ := json.MarshalIndent(data, "", "\t")

	now := time.Now().Unix()

	ioutil.WriteFile(fmt.Sprintf("%d.json", now), d, 0777)
	//

	var order model.ServiceOrderRequest

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	ord, err := env.Services.orderService.CreateServiceOrder(r.Context(), order, clientID, transactionID)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, ord, http.StatusOK)
}

// Update Order Status godoc
// @Tags Update Order
// @Summary Update Provision Order Status
// @Accept  json
// @Produce json
// @Param x-client-id header string true "client id"
// @Param req body model.UpdateOrderRequest true "Order"
// @Param id path string true "id"
// @Success 201 {object} model.Order
// @Failure 400 Incorrect body format
// @Failure 404 Order Not Found
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/{id} [put]
func updateOrder(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	var request model.UpdateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	ord, err := env.Services.orderService.UpdateOrderStatus(r.Context(), request, id, clientID)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = env.Validators.orderRequestValidator.ValidateStatus(request.Status)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if ord == nil {
		pkgHttp.ErrorResponse(w, errors.New("order not found"), http.StatusNotFound)
		return
	}

	pkgHttp.JsonResponse(w, ord, http.StatusCreated)
}
