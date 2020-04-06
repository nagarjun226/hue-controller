package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nagarjun226/hue-controller/config"
)

// HueLightT - struct for Hue Lights Data
type HueLightT struct {
	State StateT `json:"state"`
	ID    string `json:"light_id"`
}

// StateT - State struct as defined in the Hue API
type StateT struct {
	On         bool  `json:"on"`
	Brightness uint8 `json:"bri"`
}

// UpdateDevicesConnection - obtain all the devices associated with the bridge
// Makes a call to the Hue Birde `Lights` API.
func UpdateDevicesConnection(conn config.HueConnectionT) ([]HueLightT, error) {

	url, err := conn.GetEP()
	if err != nil {
		return []HueLightT{}, err
	}
	url = url + "/lights"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []HueLightT{}, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []HueLightT{}, err
	}

	defer resp.Body.Close()

	lightsMap := make(map[string]HueLightT)
	err = json.NewDecoder(resp.Body).Decode(&lightsMap)
	if err != nil {
		return []HueLightT{}, err
	}

	lights := make([]HueLightT, 0)
	for k, v := range lightsMap {
		v.ID = k
		lights = append(lights, v)
	}

	return lights, nil
}

// SetLightState - Set the state of (conn, light_id) with the parameters present in state
// makes a PUT request to the HueBridge API
func SetLightState(conn config.HueConnectionT, lightID string, state StateT) error {
	url, err := conn.GetEP()
	if err != nil {
		return err
	}
	url = url + "/lights/" + lightID + "/state"

	//fmt.Println(url)

	stateJSON, err := json.Marshal(&state)
	if err != nil {
		return err
	}

	// initialize http client
	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err1 := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(stateJSON))
	if err1 != nil {
		return err1
	}

	//fmt.Printf("REQUES == %+v\n", req)

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err2 := client.Do(req)
	if err2 != nil {
		return err2
	}

	//fmt.Printf("Resp == %+v\n", resp)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unsucessful in sending the request to the Hue Bridge. Status Code  = %v", resp.StatusCode)
	}

	return nil
}
