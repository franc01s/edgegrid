package main

type Power struct {
	SiteCurrentPowerFlow struct {
		UpdateRefreshRate int    `json:"updateRefreshRate"`
		Unit              string `json:"unit"`
		Connections       []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"connections"`
		GRID struct {
			Status       string  `json:"status"`
			CurrentPower float64 `json:"currentPower"`
		} `json:"GRID"`
		LOAD struct {
			Status       string  `json:"status"`
			CurrentPower float64 `json:"currentPower"`
		} `json:"LOAD"`
		PV struct {
			Status       string  `json:"status"`
			CurrentPower float64 `json:"currentPower"`
		} `json:"PV"`
		STORAGE struct {
			Status       string  `json:"status"`
			CurrentPower float64 `json:"currentPower"`
			ChargeLevel  int     `json:"chargeLevel"`
			Critical     bool    `json:"critical"`
		} `json:"STORAGE"`
	} `json:"siteCurrentPowerFlow"`
}
