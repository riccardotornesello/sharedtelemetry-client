package iracing

import (
	"example/sharedtelemetry/common"

	"github.com/riccardotornesello/irsdk-go"
)

type IRacingConnection struct {
	irsdk *irsdk.IRSDK
}

func NewConnection() *IRacingConnection {
	irsdk := irsdk.Init(nil)

	return &IRacingConnection{
		irsdk: irsdk,
	}
}

func (c *IRacingConnection) GetData() ([]common.Driver, common.Event) {
	c.irsdk.Update(true)

	session := c.irsdk.Session

	drivers := make([]common.Driver, len(session.DriverInfo.Drivers))

	positions := c.irsdk.Telemetry["CarIdxPosition"].Array().([]int)
	lastLapTimes := c.irsdk.Telemetry["CarIdxLastLapTime"].Array().([]float32)
	bestLapTimes := c.irsdk.Telemetry["CarIdxBestLapTime"].Array().([]float32)

	for i, driver := range session.DriverInfo.Drivers {
		carIdx := driver.CarIdx
		drivers[i] = common.Driver{
			DriverName:  driver.UserName,
			TeamName:    driver.TeamName,
			Position:    positions[carIdx],
			CarNumber:   driver.CarNumberRaw,
			LastLapTime: lastLapTimes[carIdx],
			BestLapTime: bestLapTimes[carIdx],
			Rating:      driver.IRating,
		}
	}

	event := common.Event{
		TrackName: session.WeekendInfo.TrackDisplayName,
	}

	return drivers, event
}
