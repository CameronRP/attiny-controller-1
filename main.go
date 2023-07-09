/*
attiny-controller - Communicates with ATtiny microcontroller
Copyright (C) 2018, The Cacophony Project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/TheCacophonyProject/go-config"
	arg "github.com/alexflint/go-arg"
	"periph.io/x/conn/v3/i2c"
)

// How long to wait before checking the recording window. This
// gives time to do something with the device before it turns off.
const (
	initialGracePeriod      = 20 * time.Minute
	batteryCSVFile          = "/var/log/battery.csv"
	batteryReadingInterval  = 10 * time.Minute
	systemStatFile          = "/proc/stat"
	saltCommandWaitDuration = 30 * time.Minute
)

var (
	version = "<not set>"

	mu                 sync.Mutex
	stayOnUntil        = time.Now()
	saltCommandWaitEnd = time.Time{}
)

type Args struct {
	ConfigDir          string `arg:"-c,--config" help:"configuration folder"`
	SkipWait           bool   `arg:"-s,--skip-wait" help:"will not wait for the date to update"`
	Timestamps         bool   `arg:"-t,--timestamps" help:"include timestamps in log output"`
	SkipSystemShutdown bool   `arg:"--skip-system-shutdown" help:"don't shut down operating system when powering down"`
}

func (Args) Version() string {
	return version
}

func procArgs() Args {
	args := Args{
		ConfigDir: config.DefaultConfigDir,
	}
	arg.MustParse(&args)
	return args
}

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
	}
	// If no error then keep the background goroutines running.
	//runtime.Goexit()
}

func runMain() error {
	args := procArgs()

	if !args.Timestamps {
		log.SetFlags(0)
	}

	log.Printf("running version: %s", version)

	conf, err := ParseConfig(args.ConfigDir)
	if err != nil {
		return err
	}

	log.Println("Connecting to ATtiny1616")
	attiny, err := connectToATtinyWithRetries(0)
	if err != nil {
		return err
	}

	log.Println("Setting up WDT pinging")
	go func() {
		for {
			attiny.PingWatchdog()
			time.Sleep(time.Second * 5)
		}
	}()

	log.Println("Connecting to RTC")
	rtc, err := InitPCF9564()
	if err != nil {
		return err
	}

	if err := attiny.ReadCameraState(); err != nil {
		return err
	}

	if err := rtc.ClearAlarmFlag(); err != nil {
		return err
	}

	// TODO write time to RTC if have synced with internet time thing
	//if err := rtc.SetTime(time.Now()); err != nil {
	//	return err
	//}

	t, err := rtc.GetTime()
	if err != nil {
		return err
	}
	log.Println("RTC time:", t.Format(time.RFC3339))

	alarmTime := conf.OnWindow.NextStart()
	log.Println("Alarm time:", alarmTime.Format(time.RFC3339))

	if err := rtc.SetAlarmTime(AlarmTimeFromTime(alarmTime)); err != nil {
		return err
	}

	if err := rtc.SetAlarmEnabled(true); err != nil {
		return err
	}

	attiny.ReadCameraState()
	log.Println(attiny.CameraState)

	if args.SkipWait {
		log.Println("Not waiting initial grace period.")
	} else {
		log.Println("Waiting 20 minutes before checking for power off.")
		time.Sleep(20 * time.Minute)
	}

	// Wait for next power off time if in active window or if window is going to starts in the next 5 minutes
	if conf.OnWindow.Active() || time.Until(conf.OnWindow.NextStart()) < time.Minute*2 {
		delayDuration := time.Until(conf.OnWindow.NextEnd())
		log.Printf("Sleeping for %v until turning off", delayDuration)
		time.Sleep(delayDuration)
		log.Println("Finished waiting for power off")
		attiny.ReadCameraState()
		log.Println(attiny.CameraState)
	}

	attiny.ReadCameraState()
	log.Println(attiny.CameraState)
	if err := attiny.PoweringOff(); err != nil {
		return err
	}
	attiny.ReadCameraState()
	log.Println(attiny.CameraState)
	shutdown()
	return nil
}

func fromBCD(b byte) int {
	return int(b&0x0F) + int(b>>4)*10
}

// readBytes reads bytes from the I2C device starting from a given register.
func readBytes(dev *i2c.Dev, register byte, data []byte) error {
	return dev.Tx([]byte{register}, data)
}

// readByte reads a byte from the I2C device from a given register.
func readByte(dev *i2c.Dev, register byte) (byte, error) {
	data := make([]byte, 1)
	if err := dev.Tx([]byte{register}, data); err != nil {
		return 0, err
	}
	return data[0], nil
}

// writeByte writes a byte to the I2C device at a given register.
func writeByte(dev *i2c.Dev, register byte, data byte) error {
	_, err := dev.Write([]byte{register, data})
	return err
}

// toBCD converts a decimal number to binary-coded decimal.
func toBCD(n int) byte {
	return byte(n)/10<<4 + byte(n)%10
}

// writeBytes writes the given bytes to the I2C device.
func writeBytes(dev *i2c.Dev, data []byte) error {
	_, err := dev.Write(data)
	return err
}

func shutdown() error {
	cmd := exec.Command("/sbin/poweroff")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("poweroff failed: %v\n%s", err, output)
	}
	return nil
}
