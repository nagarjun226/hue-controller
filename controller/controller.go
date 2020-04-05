package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/nagarjun226/hue-controller/config"
	"github.com/nagarjun226/hue-controller/devices"
)

// HueControllerT - A single controller that will control one Hue Bridge
type HueControllerT struct {
	Connection config.HueConnectionT
	Lights     []devices.HueLightT
	Registered bool   // Is this connection registered by the device registry service
	HumanName  string // Human name of this HueBridge
}

// HueDeviceRegistryT - struct onto which some extra device information from the device registry will be unmarshalled
type HueDeviceRegistryT struct {
	Data RegistryDataT `json:"data"`
}

// RegistryDataT - Data object of for the registry service
type RegistryDataT struct {
	HumanName string `json:"device_name"`
	Gateway   string `json:"controller_gateway"`
}

// GetHueController - Constructor for a controller given the connection
func GetHueController(conn config.HueConnectionT) (*HueControllerT, error) {

	// ToDO: Think about this? can this ever happen?
	if conn.Username == "" {
		return nil, fmt.Errorf("No username set")
	}

	if conn.Bridge == "" {
		return nil, fmt.Errorf("ipAddress of Bridge not set")
	}

	c := new(HueControllerT)
	c.Connection = conn

	return c, nil
}

// Controllers - Global variable that keeps track of all the controllers
var Controllers struct {
	Cons []*HueControllerT
	sync.Mutex
	ConSet bool // Let other packages know if the controllers were set
}

// CreateControllersFromConfig - Look through the config and create a separate controller for each bridge connection
func CreateControllersFromConfig() (errMsg error) {
	config.Config.Lock()
	defer config.Config.Unlock()

	Controllers.Lock()
	defer Controllers.Unlock()

	var failedConnections = make([]int, 0) // Keep of track of the index number of the failed connections
	for ii, conn := range config.Config.Connections {
		// Create Controller
		c, err := GetHueController(conn)
		if err != nil {
			//fmt.Println(err)
			failedConnections = append(failedConnections, ii)
			continue
		}

		err1 := c.update()
		if err1 != nil {
			//fmt.Println(err1)
			failedConnections = append(failedConnections, ii)
			continue
		}

		Controllers.Cons = append(Controllers.Cons, c)
	}
	if len(failedConnections) > 0 {
		errMsg = fmt.Errorf("Failed connections at connections = %v", failedConnections)
	} else {
		errMsg = nil
	}
	Controllers.ConSet = true

	return
}

// Update the controller with registry details
func getHueBridgeRegistryDetails(conn config.HueConnectionT) (HueDeviceRegistryT, bool, error) {
	uri := config.REGISTRYDOMAIN + "/util/search"
	uri = fmt.Sprintf("%v?controller_gateway=%v", uri, conn.Bridge)
	client := &http.Client{}

	//data := make(map[string]interface{})
	//data["controller_gateway"] = con.Connection.Bridge
	//jsonData, e := json.Marshal(&data)
	//if e != nil {
	//	return e
	//}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return HueDeviceRegistryT{}, false, err
	}
	req.Header.Add("Accept", "application/json")
	//fmt.Printf("===%+v\n", req)

	resp, err := client.Do(req)
	if err != nil {
		return HueDeviceRegistryT{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return HueDeviceRegistryT{}, false, fmt.Errorf("connection not registered device not found conn = %+v", conn)
	}

	respS := new(HueDeviceRegistryT)
	err = json.NewDecoder(resp.Body).Decode(respS)

	if err != nil {
		return HueDeviceRegistryT{}, false, err
	}

	return *respS, true, nil
}

func init() {
	Controllers.Lock()
	Controllers.Cons = make([]*HueControllerT, 0)
	Controllers.Unlock()
	var err error
	for {
		if config.ConfigSet {
			err = CreateControllersFromConfig()
			break
		} else {
			time.Sleep(time.Duration(2) * time.Second)
		}
	}

	go updateControllersChron()

	fmt.Println(err)
}

// Method that keeps the controller details up to date with the sources of truth
// This is unsafe
// Lights - updated from the HueBridge API
// Registry details - updated from device registry
func (con *HueControllerT) update() error {
	// Update the Lights
	lights, err := devices.UpdateDevicesConnection(con.Connection)
	if err != nil {
		return err
	}
	con.Lights = lights

	regDetails, registered, e := getHueBridgeRegistryDetails(con.Connection)
	//fmt.Println("---", regDetails, registered, e)
	con.HumanName = regDetails.Data.HumanName
	con.Registered = registered

	if e != nil {
		return err
	}

	return nil
}

// Chron job that updates the controllers every 30s
func updateControllersChron() {
	for {
		time.Sleep(time.Duration(3) * time.Second)

		if !Controllers.ConSet {
			CreateControllersFromConfig()
		}

		Controllers.Lock()

		for ii := range Controllers.Cons {
			Controllers.Cons[ii].update()
			//fmt.Println(err)
			//fmt.Printf("%+v,,%v\n", *Controllers.Cons[ii], time.Now())
		}

		Controllers.Unlock()
	}
}
