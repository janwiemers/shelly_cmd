package shelly_cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type callPayload struct {
	Method string
	Path   string
}

// HttpApi provides the functionality to communicate with the shelly.
// It is used by all other functions in this package.
// It contains the IP address of the shelly and the online status.
// The IP address is used to construct the URL for the API calls.
// The online status is used to determine if the shelly is reachable.
// The client is used to make the API calls.
// The client is initialized with a default timeout of 5 seconds.
// The client is also initialized with a default transport that disables keep-alive connections.
type RpcApi struct {
	IP     string
	online bool
	client *http.Client
}

// NewHttpApi creates a new instance of the HttpApi struct.
// It takes the IP address of the shelly as a parameter.
// The IP address is used to construct the URL for the API calls.
// The client is used to make the API calls.
// The client is initialized with a default timeout of 5 seconds.
// The client is also initialized with a default transport that disables keep-alive connections.
func NewRpcApi(ip string) *RpcApi {
	return &RpcApi{
		IP:     ip,
		online: false,
		client: &http.Client{},
	}
}

// call provides the internal central point of communicating with the shelly.
// It is used by all other functions in this package.
// It returns the response body as a byte array and an error if any occurred.
// The error is returned in case of a non 2xx response code.
// The caller is responsible for parsing the response body.
//
//	func (h *HttpApi) callGet(path string) ([]byte, error) {
//		return h.call(callPayload{
//			Method: "GET",
//			Path:   path,
//		})
//	}
func (r *RpcApi) call(payload callPayload) ([]byte, error) {
	url := fmt.Sprintf("http://%s/rpc/%s", r.IP, payload.Path)

	req, err := http.NewRequest(payload.Method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *RpcApi) SwitchOn(relay int) error {
	_, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.Set?id=%s&on=true", strconv.Itoa(relay)),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *RpcApi) SwitchOff(relay int) error {
	_, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.Set?id=%s&on=false", strconv.Itoa(relay)),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *RpcApi) SwitchOnWithTimer(relay, timer int) error {
	_, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.Set?id=%s&on=true&toggle_after=%s", strconv.Itoa(relay), strconv.Itoa(timer)),
	})

	if err != nil {
		return err
	}
	return nil
}

type DeviceInfoResponse struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Mac                   string `json:"mac"`
	Model                 string `json:"model"`
	Version               string `json:"ver"`
	ForwardVersion        string `json:"fw_id"`
	AuthenticationEnabled bool   `json:"auth_en"`
	AuthenticationDomain  string `json:"auth_domain"`
	Generation            int    `json:"gen"`
}

type SwitchStatusResponse struct {
	ID           int     `json:"id"`
	Source       string  `json:"source"`
	Output       bool    `json:"output"`
	Apower       float64 `json:"apower"`
	Voltage      float64 `json:"voltage"`
	Current      float64 `json:"current"`
	AFrequency   float64 `json:"freq"`
	ActiveEnergy struct {
		Total    float64   `json:"total"`
		ByMinute []float64 `json:"by_minute"`
		MinuteTs int       `json:"minute_ts"`
	} `json:"aenergy"`
	ReturnedEnergy struct {
		Total    float64   `json:"total"`
		ByMinute []float64 `json:"by_minute"`
		MinuteTs int       `json:"minute_ts"`
	} `json:"ret_aenergy"`
	Temperature struct {
		TC float64 `json:"tC"`
		TF float64 `json:"tF"`
	} `json:"temperature"`
}

// SwitchStatus retrieves the status of a specific relay.
// It returns a SwitchStatusResponse struct and an error if any occurred.
// The error is returned in case of a non 2xx response code.
// The caller is responsible for parsing the response body.
func (r *RpcApi) SwitchStatus(relay int) (*SwitchStatusResponse, error) {
	res, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.GetStatus?id=%s", strconv.Itoa(relay)),
	})

	if err != nil {
		return nil, err
	}

	var result SwitchStatusResponse
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type SwitchResetResponse struct {
	ActiveEnergy struct {
		Total float32 `json:"total"`
	}
	ReturnEnergy struct {
		Total float32 `json:"total"`
	}
}

func (r *RpcApi) SwitchReset(relay int) (*SwitchResetResponse, error) {
	res, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.ResetCounters?id=%s&type=%s", strconv.Itoa(relay), "['aenergy','ret_aenergy']"),
	})

	if err != nil {
		return nil, err
	}

	var result SwitchResetResponse
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type SwitchConfigResponse struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Mode              string `json:"in_mode"`
	InitialState      string `json:"initial_state"`
	AutoOn            bool   `json:"auto_on"`
	AutoOnDelay       int    `json:"auto_on_delay"`
	AutoOff           bool   `json:"auto_off"`
	AutoOffDelay      int    `json:"auto_off_delay"`
	Autorecover       bool   `json:"autorecover_voltage_errors"`
	PowerLimit        int    `json:"power_limit"`
	VoltageLimit      int    `json:"voltage_limit"`
	UndervoltageLimit int    `json:"undervoltage_limit"`
	CurrentLimit      int    `json:"current_limit"`
}

func (r *RpcApi) SwitchConfig(relay int) (*SwitchConfigResponse, error) {
	res, err := r.call(callPayload{
		Method: "GET",
		Path:   fmt.Sprintf("Switch.GetConfig?id=%s", strconv.Itoa(relay)),
	})

	if err != nil {
		return nil, err
	}

	var result SwitchConfigResponse
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeviceInfo retrieves the device information from the shelly.
// It returns a DeviceInfoResponse struct and an error if any occurred.
// The error is returned in case of a non 2xx response code.
// The caller is responsible for parsing the response body.
func (r *RpcApi) GetDeviceInfo() (*DeviceInfoResponse, error) {
	res, err := r.call(callPayload{
		Method: "GET",
		Path:   "Shelly.GetDeviceInfo",
	})

	if err != nil {
		return nil, err
	}

	var result DeviceInfoResponse
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
