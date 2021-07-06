package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../repos"
	"../middleware"
	"net/http"
	"../models"
)

type TransferController struct {
	Repo repos.Repository
}

func (c *TransferController) AddRoutes(r *mux.Router) {
	rAuthPlayer := r.PathPrefix("/transfer").Subrouter()
	rAuthPlayer.Use(middleware.Auth(c.Repo))
	rAuthPlayer.HandleFunc("/all", c.handleGETAll).Methods("GET")
	rAuthPlayer.HandleFunc("/{id}", c.handleGET).Methods("GET")
	rAuthPlayer.HandleFunc("/{id}", c.handlePATCH).Methods("PATCH")
	rAuthPlayer.HandleFunc("/{id}", c.handleDELETE).Methods("DELETE")
	rAuthPlayer.HandleFunc("", c.handlePOST).Methods("POST")
}

func (c *TransferController) handleGETAll(w http.ResponseWriter, req *http.Request) {
	panic("filter by params")
	transfers := c.Repo.GetTransfers()

	arr := make([]transferJson, 0)
	for _, transfer := range transfers {
		arr = append(arr, c.makeJson(&transfer))
	}

	resp, err := json.Marshal(arr)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, resp)
}

func (c *TransferController) handleGET(w http.ResponseWriter, req *http.Request) {
	transfer, err := c.getTransferFromRequest(w, req)
	if err != nil {
		return
	}

	resp, err := json.Marshal(c.makeJson(&transfer))
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, resp)
}

// Handle editing a transfer
func (c *TransferController) handlePATCH(w http.ResponseWriter, req *http.Request) {

}

// Handle deleting a transfer
func (c *TransferController) handleDELETE(w http.ResponseWriter, req *http.Request) {
	transfer, err1 := c.getTransferFromRequest(w, req)
	user, err2 := getAuthenticatedUserFromRequest(w, req)
	if err1 != nil || err2 != nil {
		return
	}

	if !user.IsAdmin() && user.ID != transfer.Player.Team.OwnerID {
		writeUnauthorizedError(w, fmt.Errorf("trying to delete a not owned transfer"))
		return
	}

	err := c.Repo.Delete(transfer)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	writeOkResponse(w)
}

type transferJson struct {

}

// Handle creating a transfer
func (c *TransferController) handlePOST(w http.ResponseWriter, req *http.Request) {

}

// Gets transfers from the request
func (c *TransferController) getTransferFromRequest(w http.ResponseWriter, req *http.Request) (models.Transfer, error) {
	id, err := parseIdFromRequest(w, req)
	transfer, err := c.Repo.GetTransfer(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return models.Transfer{}, err
	}
	return transfer, nil
}

// Make json from a transfer model
func (c *TransferController) makeJson(transfer *models.Transfer) transferJson {

}