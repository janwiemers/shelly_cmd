package main

import "time"

type ShellyStatusResponse struct {
	Wifi struct {
		Connected bool   `json:"connected"`
		SSID      string `json:"ssid"`
		IP        string `json:"ip"`
		RSSI      string `json:"rssi"`
	} `json:"wifi_sta"`
	Cloud struct {
		Enabled   bool `json:"enabled"`
		Connected bool `json:"connected"`
	} `json:"cloud"`
	MQTT struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	Time         string    `json:"time"`
	Unixtime     time.Time `json:"unixtime"`
	Serial       int       `json:"serial"`
	HasUpdate    bool      `json:"has_update"`
	Mac          string    `json:"mac"`
	CfgChanged   int       `json:"cfg_changed"`
	ActionsStats struct {
		Skipped int `json:"skipped"`
	} `json:"actions_stats"`
	Relays []Relay `json:"relays"`
	Meters []struct {
		Power     int       `json:"power"`
		Overpower int       `json:"overpower"`
		IsValid   bool      `json:"is_valid"`
		Timestamp time.Time `json:"timestamp"`
		Counters  []int     `json:"counters"`
		Total     int       `json:"total"`
	} `json:"meters"`
	Inputs []struct {
		Input    int    `json:"input"`
		Event    string `json:"event"`
		EventCnt int    `json:"event_cnt"`
	} `json:"inputs"`
	Temperature     float32 `json:"temperature"`
	Overtemperature bool    `json:"overtemperature"`
	Tmp             struct {
		TC      float32 `json:"tC"`
		TF      float32 `json:"tF"`
		IsValid bool    `json:"is_valid"`
	} `json:"tmp"`
	TemperatureStatus string `json:"temperature_status"`
	Update            struct {
		Status      string `json:"status"`
		HasUpdate   bool   `json:"has_update"`
		NewVersion  string `json:"new_version"`
		OldVersion  string `json:"old_version"`
		BetaVersion string `json:"beta_version"`
	} `json:"update"`
	RamTotal int     `json:"ram_total"`
	RamFree  int     `json:"ram_free"`
	FsSize   int     `json:"fs_size"`
	FsFree   int     `json:"fs_free"`
	Voltage  float32 `json:"voltage"`
	UpTime   int     `json:"uptime"`
}
