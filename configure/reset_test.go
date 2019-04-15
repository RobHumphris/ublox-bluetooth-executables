package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	ub "github.com/RobHumphris/ublox-bluetooth"
	"github.com/RobHumphris/ublox-bluetooth/serial"
)

func TestReset(t *testing.T) {
	bt, err := ub.NewUbloxBluetooth(6 * time.Second)
	if err != nil {
		log.Fatalf("NewUbloxBluetooth error %v\n", err)
	}
	defer bt.Close()

	serial.SetVerbose(true)

	err = bt.SetCommsRate(serial.Default)
	if err != nil {
		fmt.Printf("Baud Rate reset error %v\n", err)
		return
	}

	bt.EnterCommandMode()

	err = bt.MultipleATCommands()
	if err != nil {
		err = bt.SetCommsRate(serial.Default)
		if err != nil {
			fmt.Printf("Baud Rate reset error %v\n", err)
			return
		}
	}
	fmt.Println("FactoryReset AT response OK")

	err = bt.FactoryReset()
	if err != nil {
		fmt.Printf("FactoryReset error %v\n", err)
		return
	}
	fmt.Println("FactoryReset OK")

	err = bt.RebootUblox()
	if err != nil {
		fmt.Printf("FactoryReset reboot error %v\n", err)
		return
	}
	fmt.Println("FactoryReset reboot OK")

	bt.EnterCommandMode()

	err = bt.MultipleATCommands()
	if err != nil {
		log.Fatalf("NewUbloxBluetooth AT error %v\n", err)
		return
	}
	fmt.Println("FactoryReset AT - 1 response OK")
}

func TestConfigure(t *testing.T) {
	bt, err := ub.NewUbloxBluetooth(6 * time.Second)
	if err != nil {
		log.Fatalf("NewUbloxBluetooth error %v\n", err)
	}
	defer bt.Close()

	defaultSettings(bt)
}
