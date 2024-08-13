package common

type Driver struct {
	Id         int    `json:"id"`
	DriverName string `json:"driverName"`
	TeamName   string `json:"teamName"`
	CarNumber  int    `json:"carNumber"`
	Rating     int    `json:"rating"`
}

type Session struct {
	TrackName string `json:"trackName"`
}

type InputTelemetry struct {
	Throttle    float32 `json:"throttle"`
	Brake       float32 `json:"brake"`
	Steering    float32 `json:"steering"`
	SteeringMax float32 `json:"steeringMax"`
	Clutch      float32 `json:"clutch"`
}

type CarTelemetry struct {
	Tyres [4]Tyre `json:"tyres"`
}

type Tyre struct {
	TempCarcass [3]float32 `json:"tempCarcass"`
	TempSurface [3]float32 `json:"tempSurface"`
}

type Radio struct {
	SpeakingCarIdx int `json:"speakingCarIdx"`
}
