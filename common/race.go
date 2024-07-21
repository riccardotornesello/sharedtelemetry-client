package common

type Driver struct {
	DriverName  string  `json:"driverName"`
	TeamName    string  `json:"teamName"`
	Position    int     `json:"position"`
	CarNumber   int     `json:"carNumber"`
	LastLapTime float32 `json:"lastLapTime"`
	BestLapTime float32 `json:"bestLapTime"`
	Rating      int     `json:"rating"`
}

type Event struct {
	TrackName string `json:"trackName"`
}
