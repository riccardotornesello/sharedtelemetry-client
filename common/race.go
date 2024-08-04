package common

type Driver struct {
	Id          int     `json:"id"`
	DriverName  string  `json:"driverName"`
	TeamName    string  `json:"teamName"`
	Position    int     `json:"position"`
	CarNumber   int     `json:"carNumber"`
	LastLapTime float32 `json:"lastLapTime"`
	BestLapTime float32 `json:"bestLapTime"`
	Rating      int     `json:"rating"`
	Gap         float32 `json:"gap"`
}

type Race struct {
	TrackName string `json:"trackName"`
}

type Telemetry struct {
	Throttle float32 `json:"throttle"`
	Brake    float32 `json:"brake"`
	Steer    float32 `json:"steer"`

	Tyres [4]Tyre `json:"tyres"`
}

type Tyre struct {
	TempCarcass [3]float32 `json:"tempCarcass"`
	TempSurface [3]float32 `json:"tempSurface"`
}
