package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/srrmendez/private-api-offers/model"
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

// Search Offers godoc
// @Tags Search Offers
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param active query bool false "offers status"
// @Success 200 {array} model.Offer
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/ [get]
func searchOffers(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var active *bool

	if r.URL.Query().Get("active") != "" {
		st, _ := strconv.ParseBool(r.URL.Query().Get("active"))
		active = &st
	}

	offers, err := env.offerService.Search(r.Context(), clientID, active)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, offers, http.StatusOK)
}

// Get Offer godoc
// @Tags Get Offer
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param id path string true "id"
// @Success 200 {object} model.Offer
// @Failure 404 Offer Not Found
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/{id} [get]
func getOffer(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	offer, err := env.offerService.Get(r.Context(), id, clientID)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if offer == nil {
		pkgHttp.ErrorResponse(w, errors.New("order not found"), http.StatusNotFound)
		return
	}

	pkgHttp.JsonResponse(w, offer, http.StatusOK)
}

// Create Offers godoc
// @Tags Create Offers Order
// @Summary Create Offers from commercial system
// @Accept  json
// @Produce json
// @Param x-client-id header string true "client id"
// @Param req body model.BssSyncOfferRequest true "offers to sync"
// @Success 201
// @Failure 400 Incorrect body format
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/ [post]
func createOffers(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// TODO Remove this
	go func() {
		var data interface{}
		json.NewDecoder(r.Body).Decode(&data)

		d, _ := json.MarshalIndent(data, "", "\t")

		now := time.Now().Unix()

		ioutil.WriteFile(fmt.Sprintf("%d.json", now), d, 0777)
	}()
	//

	var request model.BssSyncOfferRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = env.offerService.Create(r.Context(), clientID, request)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, map[string]string{}, http.StatusCreated)
}
