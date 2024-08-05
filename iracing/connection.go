package iracing

import (
	"fmt"
	"sharedtelemetry/client/common"
	"time"

	"github.com/riccardotornesello/irsdk-go"
)

var tyreCodes = [4]string{"LF", "RF", "LR", "RR"}
var tyrePos = [3]string{"L", "M", "R"}

type IRacingConnection struct {
	irsdk *irsdk.IRSDK
	quit  chan struct{}

	isConnected bool
	drivers     *[]common.Driver
	race        *common.Race
	telemetry   *common.Telemetry

	chatTalkingIdx int
}

func NewConnection() *IRacingConnection {
	irsdk := irsdk.Init(nil)

	return &IRacingConnection{
		irsdk: irsdk,

		drivers:   &[]common.Driver{},
		race:      &common.Race{},
		telemetry: &common.Telemetry{},

		chatTalkingIdx: -1,
	}
}

func (c *IRacingConnection) Start(updateDelay int, connectionDelay int, eventChannel chan common.Event) {
	c.quit = make(chan struct{})

	go c.refetchData(c.quit, updateDelay, 3)

	for {
		start := time.Now()

		c.isConnected = c.irsdk.IsConnected()

		if c.isConnected {

			session := c.irsdk.Session

			drivers := make([]common.Driver, len(session.DriverInfo.Drivers))

			positions := c.irsdk.Telemetry["CarIdxPosition"].Array().([]int)
			lastLapTimes := c.irsdk.Telemetry["CarIdxLastLapTime"].Array().([]float32)
			bestLapTimes := c.irsdk.Telemetry["CarIdxBestLapTime"].Array().([]float32)
			gaps := c.irsdk.Telemetry["CarIdxEstTime"].Array().([]float32)

			throttle := c.irsdk.Telemetry["Throttle"].Value().(float32)
			brake := c.irsdk.Telemetry["Brake"].Value().(float32)
			steerAngle := c.irsdk.Telemetry["SteeringWheelAngle"].Value().(float32)
			steerMax := c.irsdk.Telemetry["SteeringWheelAngleMax"].Value().(float32)
			steer := steerAngle / steerMax

			tyres := [4]common.Tyre{}
			for i, code := range tyreCodes {
				tyres[i] = common.Tyre{
					TempCarcass: [3]float32{},
					TempSurface: [3]float32{},
				}
				for j, pos := range tyrePos {
					carcass, ok := c.irsdk.GetVar(code + "tempC" + pos)
					if !ok {
						tyres[i].TempCarcass[j] = 0
					} else {
						tyres[i].TempCarcass[j] = carcass.(float32)
					}

					surface, ok := c.irsdk.GetVar(code + "temp" + pos)
					if !ok {
						tyres[i].TempSurface[j] = 0
					} else {
						tyres[i].TempSurface[j] = surface.(float32)
					}
				}
			}

			for i, driver := range session.DriverInfo.Drivers {
				carIdx := driver.CarIdx
				drivers[i] = common.Driver{
					Id:          driver.UserID,
					DriverName:  driver.UserName,
					TeamName:    driver.TeamName,
					Position:    positions[carIdx],
					CarNumber:   driver.CarNumberRaw,
					LastLapTime: lastLapTimes[carIdx],
					BestLapTime: bestLapTimes[carIdx],
					Rating:      driver.IRating,
					Gap:         gaps[carIdx],
				}
			}

			race := common.Race{
				TrackName: session.WeekendInfo.TrackDisplayName,
			}

			telemetry := common.Telemetry{
				Throttle: throttle,
				Brake:    brake,
				Steer:    steer,

				Tyres: tyres,
			}

			c.drivers = &drivers
			c.race = &race
			c.telemetry = &telemetry

			elapsed := time.Since(start)
			time.Sleep(time.Duration(updateDelay)*time.Millisecond - elapsed)
		} else {
			time.Sleep(time.Duration(connectionDelay) * time.Millisecond)
		}
	}
}

func (c *IRacingConnection) Stop() {
	close(c.quit)
}

func (c *IRacingConnection) GetData() (*[]common.Driver, *common.Race, *common.Telemetry) {
	return c.drivers, c.race, c.telemetry
}

func (c *IRacingConnection) refetchData(quit chan struct{}, refreshRate int, sessionSkip int) {
	i := 0

	for {
		select {
		case <-quit:
			return
		default:
			start := time.Now()

			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Recovered from panic: %v\n", r)
					}
				}()

				c.irsdk.Update(i == 0)
				c.isConnected = c.irsdk.IsConnected()

			}()

			i = (i + 1) % sessionSkip
			time.Sleep(time.Second/time.Duration(refreshRate) - time.Since(start))
		}
	}
}
