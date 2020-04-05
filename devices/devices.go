package devices

import (
	"encoding/json"
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
