package controllers

import (
	"github.com/gorilla/mux"
	"../repos"
)

type TransferController struct {
	Repo repos.Repository
}

func (c *TransferController) AddRoutes(r *mux.Router) {

}