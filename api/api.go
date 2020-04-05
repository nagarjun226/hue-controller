package api

import (
	"github.com/gorilla/mux"
	"github.com/nagarjun226/hue-controller/controller"
)

// API - API instance
type API struct {
}

// GetAPIInstance - Return an Instance of the API
func GetAPIInstance() *API {
	a := new(API)
	return a
}

// Router - Router for the API
func (api *API) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/controllers", controller.ListControllersEP).Methods("GET")
	r.HandleFunc("/api/{bridge}/lights", controller.GetBridgeLights).Methods("GET")
	r.HandleFunc("/api/{bridge}/lights/{light_id}", controller.GetBridgeLight).Methods("GET")
	return r
}
