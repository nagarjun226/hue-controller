package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// SERVICENAME - name of the service
const SERVICENAME = "hue-controller"

// CONFIGDOMAIN - Domain on which config service is running
const CONFIGDOMAIN = "http://localhost:8080"

// REGISTRYDOMAIN - Domain on which the device registry service is runnning
const REGISTRYDOMAIN = "http://localhost:5000"

// GETCONFEP - getconfig endpoint for our service
const GETCONFEP = "/getconfig/" + SERVICENAME

// ConfigSet - Let other packages know if the config has been successfully set
var ConfigSet bool

// HueConfigT - Config structure for the hue service
type HueConfigT struct {
	Connections []HueConnectionT `json:"connections"`
	sync.Mutex
}

// HueConnectionT - Struct that holds the details of a connection.
type HueConnectionT struct {
	Bridge   string `json:"bridge"`
	Username string `json:"username"`
}

// Config - Pointer to the Config
var Config *HueConfigT

// GetEP - Get the endpoint to the hue bridge from the connection
func (h HueConnectionT) GetEP() (string, error) {
	if h.Bridge == "" {
		return "", fmt.Errorf("Bridge not set for connection")
	}

	if h.Username == "" {
		return "", fmt.Errorf("Username not set for connection")
	}

	return fmt.Sprintf("http://%s/api/%s", h.Bridge, h.Username), nil
}

// SetLocalConfig - Request the configmanager for the config of hue controller
func SetLocalConfig() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", CONFIGDOMAIN+GETCONFEP, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Deconde the JSON response into the local config object
	Config.Lock()
	err = json.NewDecoder(resp.Body).Decode(&Config)
	Config.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Config = new(HueConfigT)
	err := SetLocalConfig()
	if err != nil {
		panic(err)
	}
	ConfigSet = true
}
