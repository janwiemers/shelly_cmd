package shelly_cmd

import (
	"context"
	"fmt"

	"github.com/ybbus/jsonrpc/v3"
)

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
	client jsonrpc.RPCClient
}

// NewHttpApi creates a new instance of the HttpApi struct.
// It takes the IP address of the shelly as a parameter.
// The IP address is used to construct the URL for the API calls.
// The client is used to make the API calls.
// The client is initialized with a default timeout of 5 seconds.
// The client is also initialized with a default transport that disables keep-alive connections.
func NewRpcApi(ip string) (*RpcApi, error) {
	rpcClient := jsonrpc.NewClientWithOpts(fmt.Sprintf("http://%s:%d/rpc", ip, 80), &jsonrpc.RPCClientOpts{
		AllowUnknownFields: true,
	})

	return &RpcApi{
		IP:     ip,
		client: rpcClient,
	}, nil
}

type SwitchWasResponse struct {
	WasOn bool `json:"was_on"`
}

func (r *RpcApi) SwitchOn(relay int) (*SwitchWasResponse, error) {
	var result SwitchWasResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.Set", map[string]interface{}{
		"id": relay,
		"on": true,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RpcApi) SwitchOnWithTimer(relay, timer int) (*SwitchWasResponse, error) {
	var result SwitchWasResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.Set", map[string]interface{}{
		"id":           relay,
		"on":           true,
		"toggle_after": timer,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RpcApi) SwitchOff(relay int) (*SwitchWasResponse, error) {
	var result SwitchWasResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.Set", map[string]interface{}{
		"id": relay,
		"on": false,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RpcApi) SwitchToggle(relay int) (*SwitchWasResponse, error) {
	var result SwitchWasResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.Toggle", map[string]interface{}{
		"id": relay,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type SwitchGetConfigResponse struct {
	ID                       int     `json:"id"`
	Name                     string  `json:"name"`
	Mode                     string  `json:"mode"`
	InitialState             string  `json:"initial_state"`
	AutoOn                   bool    `json:"auto_on"`
	AutoOnDelay              float32 `json:"auto_on_delay"`
	AutoOff                  bool    `json:"auto_off"`
	AutoOffDelay             float32 `json:"auto_off_delay"`
	AutorecoverVoltageErrors bool    `json:"autorecover_voltage_errors"`
	PowerLimit               float32 `json:"power_limit"`
	CurrentLimit             float32 `json:"current_limit"`
	UndervoltageLimit        float32 `json:"undervoltage_limit"`
	VoltageLimit             float32 `json:"voltage_limit"`
}

func (r *RpcApi) SwitchGetConfig(relay int) (*SwitchGetConfigResponse, error) {
	var result SwitchGetConfigResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.GetConfig", map[string]interface{}{
		"id": relay,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type SwitchGetStatusResponse struct {
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

// // SwitchGetStatus retrieves the status of a specific relay.
// // It returns a SwitchStatusResponse struct and an error if any occurred.
// // The error is returned in case of a non 2xx response code.
// // The caller is responsible for parsing the response body.
func (r *RpcApi) SwitchGetStatus(relay int) (*SwitchGetStatusResponse, error) {
	var result SwitchGetStatusResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.GetStatus", map[string]interface{}{
		"id": relay,
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type SwitchResetCountersResponse struct {
	ActiveEnergy struct {
		Total float32 `json:"total"`
	}
	ReturnEnergy struct {
		Total float32 `json:"total"`
	}
}

func (r *RpcApi) SwitchResetCounters(relay int) (*SwitchResetCountersResponse, error) {
	var result SwitchResetCountersResponse
	err := r.client.CallFor(context.Background(), &result, "Switch.ResetCounters", map[string]interface{}{
		"id": relay,
		// "type": []string{"aenergy", "ret_aenergy"},
		"type": "['aenergy', 'ret_aenergy']",
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
