package iracing

import (
	"example/sharedtelemetry/common"
	"time"

	"github.com/riccardotornesello/irsdk-go"
)

type IRacingConnection struct {
	irsdk *irsdk.IRSDK

	drivers   *[]common.Driver
	event     *common.Event
	telemetry *common.Telemetry
}

func NewConnection() *IRacingConnection {
	irsdk := irsdk.Init(nil)

	return &IRacingConnection{
		irsdk: irsdk,

		drivers:   &[]common.Driver{},
		event:     &common.Event{},
		telemetry: &common.Telemetry{},
	}
}

func (c *IRacingConnection) Start(updateDelay int) {
	for {
		c.irsdk.Update(true)

		session := c.irsdk.Session

		drivers := make([]common.Driver, len(session.DriverInfo.Drivers))

		positions := c.irsdk.Telemetry["CarIdxPosition"].Array().([]int)
		lastLapTimes := c.irsdk.Telemetry["CarIdxLastLapTime"].Array().([]float32)
		bestLapTimes := c.irsdk.Telemetry["CarIdxBestLapTime"].Array().([]float32)

		throttle := c.irsdk.Telemetry["Throttle"].Value().(float32)

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

		telemetry := common.Telemetry{
			Throttle: throttle,
		}

		c.drivers = &drivers
		c.event = &event
		c.telemetry = &telemetry

		time.Sleep(time.Duration(updateDelay) * time.Millisecond)
	}
}

func (c *IRacingConnection) GetData() (*[]common.Driver, *common.Event, *common.Telemetry) {
	return c.drivers, c.event, c.telemetry
}
