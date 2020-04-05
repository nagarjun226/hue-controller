package controller

// Put all the the controller related API handlers here

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nagarjun226/hue-controller/devices"

	"github.com/gorilla/mux"
)

// ListControllersEP - Handler function for the `GET /api/controllers`
func ListControllersEP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	Controllers.Lock()
	defer Controllers.Unlock()
	rsp, err := json.Marshal(&Controllers.Cons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(rsp))
}

// GetBridgeLights - handler function for `GET /api/{bridge}/lights` where bridge is the Human name of the philips hue Bridge
func GetBridgeLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	bridge, ok := vars["bridge"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
		return
	}

	Controllers.Lock()
	defer Controllers.Unlock()
	var present bool
	var lights []devices.HueLightT
	for _, con := range Controllers.Cons {
		if con.HumanName == bridge && con.Registered {
			lights = con.Lights
			present = true
			break
		}
	}

	if !present {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
		return
	}

	rsp, err := json.Marshal(&lights)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(rsp))

}

// GetBridgeLight - handler function for `GET /api/{bridge}/lights/{light_id}` where bridge is the Human name of the philips hue Bridge and light_id is the id of the light
func GetBridgeLight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	bridge, ok1 := vars["bridge"]
	id, ok2 := vars["light_id"]
	if !ok1 || !ok2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
	}

	Controllers.Lock()
	defer Controllers.Unlock()
	var presentController, presentLight bool
	var light devices.HueLightT
	for _, con := range Controllers.Cons {
		if con.HumanName == bridge && con.Registered {
			presentController = true
			for _, l := range con.Lights {
				if l.ID == id {
					presentLight = true
					light = l
					break
				}
			}
			break
		}
	}

	if !presentController {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Requested controller not present or registered")
		return
	}

	if !presentLight {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Requested Light not present")
		return
	}

	rsp, err := json.Marshal(&light)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("%v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(rsp))
}
