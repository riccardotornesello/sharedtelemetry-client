package iracing

import (
	"fmt"
	"reflect"
	"sharedtelemetry/client/common"
	"time"

	"github.com/riccardotornesello/irsdk-go"
)

type IRacingConnection struct {
	irsdk *irsdk.IRSDK
	quit  chan struct{}

	isConnected bool
	drivers     *[]common.Driver
}

func NewConnection() *IRacingConnection {
	irsdk := irsdk.Init(nil)

	return &IRacingConnection{
		irsdk: irsdk,
	}
}

func (c *IRacingConnection) Start(
	eventChannel chan common.Event,
	refreshRate int,
	sessionRefreshRate int,
) {
	c.quit = make(chan struct{})

	go c.refetchData(c.quit, refreshRate, sessionRefreshRate)
	go c.checkUpdates(c.quit, eventChannel, refreshRate)
}

func (c *IRacingConnection) Stop() {
	close(c.quit)
}

func (c *IRacingConnection) checkUpdates(
	quit chan struct{},
	eventChannel chan common.Event,
	refreshRate int,
) {
	session := c.irsdk.Session

	drivers := fetchDrivers(session)
	sessionInfo := fetchSessionInfo(session)
	radioSpeakingCarIdx := -1
	flags := make([]common.Flag, 0)

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

				if c.irsdk.Session != session {
					session = c.irsdk.Session

					newDrivers := fetchDrivers(session)
					if !reflect.DeepEqual(drivers, newDrivers) {
						drivers = newDrivers
						eventChannel <- common.Event{
							Event: common.EventDrivers,
							Data:  drivers,
						}
					}

					newSessionInfo := fetchSessionInfo(session)
					if !reflect.DeepEqual(sessionInfo, newSessionInfo) {
						sessionInfo = newSessionInfo
						eventChannel <- common.Event{
							Event: common.EventSession,
							Data:  sessionInfo,
						}
					}
				}

				inputTelemetry := fetchInputTelemetry(c.irsdk)
				eventChannel <- common.Event{
					Event: common.EventInputTelemetry,
					Data:  inputTelemetry,
				}

				carTelemetry := fetchCarTelemetry(c.irsdk)
				eventChannel <- common.Event{
					Event: common.EventCarTelemetry,
					Data:  carTelemetry,
				}

				newRadioSpeakingCarIdx, ok := c.irsdk.GetVar("RadioTransmitCarIdx")
				if !ok {
					newRadioSpeakingCarIdx = -1
				}
				if newRadioSpeakingCarIdx != radioSpeakingCarIdx {
					radioSpeakingCarIdx = newRadioSpeakingCarIdx.(int)
					eventChannel <- common.Event{
						Event: common.EventRadio,
						Data: common.Radio{
							SpeakingCarIdx: radioSpeakingCarIdx,
						},
					}
				}

				newFlags := make([]common.Flag, 0)
				newFlagsValue, ok := c.irsdk.GetVar("SessionFlags")
				if ok {
					newFlags = fetchFlags(newFlagsValue.(int))
				}
				if !reflect.DeepEqual(flags, newFlags) {
					flags = newFlags
					eventChannel <- common.Event{
						Event: common.EventFlags,
						Data:  flags,
					}
				}
			}()

			time.Sleep(time.Second/time.Duration(refreshRate) - time.Since(start))
		}
	}
}

func (c *IRacingConnection) refetchData(quit chan struct{}, refreshRate int, sessionRefreshRate int) {
	sessionInterval := max(1, refreshRate/sessionRefreshRate)
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

			i = (i + 1) % sessionInterval
			time.Sleep(time.Second/time.Duration(refreshRate) - time.Since(start))
		}
	}
}
