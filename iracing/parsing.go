package iracing

import (
	"sharedtelemetry/client/common"

	"github.com/riccardotornesello/irsdk-go"
)

func fetchDrivers(session *irsdk.Session) *[]common.Driver {
	drivers := make([]common.Driver, 64)

	for _, driver := range session.DriverInfo.Drivers {
		drivers[driver.CarIdx] = common.Driver{
			Id:         driver.UserID,
			DriverName: driver.UserName,
			TeamName:   driver.TeamName,
			CarNumber:  driver.CarNumberRaw,
			Rating:     driver.IRating,
		}
	}

	return &drivers
}

func fetchSessionInfo(session *irsdk.Session) *common.Session {
	return &common.Session{
		TrackName: session.WeekendInfo.TrackDisplayName,
	}
}

func fetchInputTelemetry(irsdk *irsdk.IRSDK) *common.InputTelemetry {
	var throttleValue float32
	var brakeValue float32
	var steeringValue float32
	var steeringMaxValue float32
	var clutchValue float32

	throttle, ok := irsdk.GetVar("Throttle")
	if ok {
		throttleValue = throttle.(float32)
	}

	brake, ok := irsdk.GetVar("Brake")
	if ok {
		brakeValue = brake.(float32)
	}

	steering, ok := irsdk.GetVar("SteeringWheelAngle")
	if ok {
		steeringValue = steering.(float32)
	}

	steeringMax, ok := irsdk.GetVar("SteeringWheelAngleMax")
	if ok {
		steeringMaxValue = steeringMax.(float32)
	}

	clutch, ok := irsdk.GetVar("Clutch")
	if ok {
		clutchValue = clutch.(float32)
	}

	return &common.InputTelemetry{
		Throttle:    throttleValue,
		Brake:       brakeValue,
		Steering:    steeringValue,
		SteeringMax: steeringMaxValue,
		Clutch:      clutchValue,
	}
}

func fetchCarTelemetry(irsdk *irsdk.IRSDK) *common.CarTelemetry {
	tyres := [4]common.Tyre{}
	for i, code := range tyreCodes {
		tyres[i] = common.Tyre{
			TempCarcass: [3]float32{},
			TempSurface: [3]float32{},
		}
		for j, pos := range tyrePos {
			carcass, ok := irsdk.GetVar(code + "tempC" + pos)
			if !ok {
				tyres[i].TempCarcass[j] = 0
			} else {
				tyres[i].TempCarcass[j] = carcass.(float32)
			}

			surface, ok := irsdk.GetVar(code + "temp" + pos)
			if !ok {
				tyres[i].TempSurface[j] = 0
			} else {
				tyres[i].TempSurface[j] = surface.(float32)
			}
		}
	}

	return &common.CarTelemetry{
		Tyres: tyres,
	}
}
