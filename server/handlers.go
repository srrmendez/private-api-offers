package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
// @Param category query string false "offers categories"
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

	if act := r.URL.Query().Get("active"); act != "" {
		st, _ := strconv.ParseBool(act)
		active = &st
	}

	var category *model.CategoryType

	if cat := r.URL.Query().Get("category"); cat != "" {
		err := checkRequestCategoryType(cat)
		if err != nil {
			pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		st := model.CategoryType(cat)

		category = &st
	}

	offers, err := env.offerService.Search(r.Context(), clientID, active, category)
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

// Get Offer godoc
// @Tags Get Offer
// @Accept  json
// @Produce  json
// @Param x-client-id header string true "client id"
// @Param ids query string true "ids"
// @Success 200 {array} model.Offer
// @Failure 404 Offer Not Found
// @Failure 401 Unauthorized Request
// @Failure 500 Server Error
// @Router /v1/secondary[get]
func getSecondaryOffers(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client-id")

	if clientID == "" {
		pkgHttp.ErrorResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	ids := strings.Split(r.URL.Query().Get("ids"), ",")
	if len(ids) == 0 {
		pkgHttp.ErrorResponse(w, errors.New("missing ids"), http.StatusBadRequest)

		return
	}

	offers, err := env.offerService.GetSecondaryOffers(r.Context(), ids)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, offers, http.StatusOK)
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

	var request model.BssSyncOfferRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := env.offerService.Sync(r.Context(), clientID, request); err != nil {
		pkgHttp.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pkgHttp.JsonResponse(w, map[string]string{}, http.StatusCreated)
}

func checkRequestCategoryType(cat string) error {
	categories := []model.CategoryType{model.CategoryTypeDataCenter, model.CategoryTypeYellowPages}

	for _, category := range categories {
		if cat == string(category) {
			return nil
		}
	}

	return fmt.Errorf("incorrect category posible values are %s, %s", model.CategoryTypeDataCenter, model.CategoryTypeYellowPages)
}
