package main

type Relay struct {
	IsOn            bool   `json:"ison"`
	HasTimer        bool   `json:"has_timer"`
	TimerStarted    int    `json:"timer_started"`
	TimerDuration   int    `json:"timer_duration"`
	TimerRemaining  int    `json:"timer_remaining"`
	Source          string `json:"source"`
	Overpower       bool   `json:"overpower"`
	Overtemperature bool   `json:"overtemperature"`
	IsValid         bool   `json:"is_valid"`
}
