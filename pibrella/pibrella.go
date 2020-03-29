package pibrella

import (
	"fmt"
	"time"

	"github.com/matryer/runner"
	"github.com/stianeikeland/go-rpio"
)

var (
	Green, Yellow, Red Light

	blinkTime = defaultBlinkTime

	// stop tasks for each pin
	tasks map[rpio.Pin]*runner.Task
)

const (
	// Pibrella pins, these are BCM
	PB_PIN_LIGHT_RED    = 27
	PB_PIN_LIGHT_YELLOW = 17
	PB_PIN_LIGHT_GREEN  = 4

	// onboard buzzer
	PB_PIN_BUZZER = 18

	defaultBlinkTime = 500 * time.Millisecond
)

func light(pinNr int) Light {
	pin := rpio.Pin(pinNr)
	pin.Output()
	pin.Low()
	return Light{pin}
}

// Light is an LED on the Pibrella board
type Light struct {
	Pin rpio.Pin
}

func (l Light) On() {
	l.Stop()
	l.Pin.High()
}

func (l Light) Off() {
	l.Stop()
	l.Pin.Low()
}

func (l Light) Toggle() {
	l.Stop()
	l.Pin.Toggle()
}

func (l Light) Blink() {

	l.Stop()

	// start new blinking task
	tasks[l.Pin] = runner.Go(func(shouldStop runner.S) error {
		for {
			fmt.Printf("toggle pin %d\n", l.Pin)
			l.Pin.Toggle()
			if shouldStop() {
				break
			}
			time.Sleep(blinkTime)
		}
		l.Pin.Low()
		return nil // no errors
	})

}

func (l Light) Stop() {
	// stop previous task if exists
	if task, ok := tasks[l.Pin]; ok {
		task.Stop()
		select {
		case <-task.StopChan():
			// task successfully stopped
		case <-time.After(1 * time.Second):
			// task didn't stop in time
		}
		l.Pin.Low()
	}
}

func Open() error {
	if err := rpio.Open(); err != nil {
		return err
	}

	// setup lights
	Green = light(PB_PIN_LIGHT_GREEN)
	Yellow = light(PB_PIN_LIGHT_YELLOW)
	Red = light(PB_PIN_LIGHT_RED)

	// init stop task map
	tasks = make(map[rpio.Pin]*runner.Task)

	// TODO buzzer (with software PWM)
	// TODO input?

	return nil
}

func Close() error {
	return rpio.Close()
}
