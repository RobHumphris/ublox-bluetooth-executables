package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	ub "github.com/8power/ublox-bluetooth"
	"github.com/RobHumphris/veh-structs/vehdata"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

var password = []byte{'A', 'B', 'C'}

func main() {
	mac := flag.String("mac", "", "MAC of the device")
	slotNumber := flag.Uint("slot", 0, "Slot number to download")
	flag.Parse()

	fmt.Printf("Opening device MAC %s\n", *mac)

	bt, err := ub.NewUbloxBluetooth(6 * time.Second)
	if err != nil {
		log.Fatalf("NewUbloxBluetooth error %v\n", err)
	}
	defer bt.Close()

	bt.ConnectToDevice(*mac, func() error {
		defer bt.DisconnectFromDevice()

		err := bt.EnableIndications()
		if err != nil {
			return err
		}

		err = bt.EnableNotifications()
		if err != nil {
			return err
		}

		unlocked, err := bt.UnlockDevice(password)
		if err != nil {
			return err
		}
		if !unlocked {
			return fmt.Errorf("Unlock device failed")
		}

		info, err := bt.GetInfo()
		if err != nil {
			return err
		}
		fmt.Printf("[GetInfo] replied with: %v\n", info)

		slot := int(*slotNumber)

		slotInfo, err := bt.ReadSlotInfo(slot)
		if err != nil {
			return errors.Wrap(err, "ReadSlotInfo error")
		}

		vehSlotInfo := &vehdata.VehSlotInfo{
			Time:           uint32(slotInfo.Time),
			Crc:            uint32(slotInfo.Slot),
			DWords:         uint32(slotInfo.Bytes),
			Odr:            slotInfo.SampleRate,
			Temperature:    uint32(slotInfo.Temperature),
			BatteryVoltage: uint32(slotInfo.BatteryVoltage),
			VoltageIn:      uint32(slotInfo.VoltageIn),
		}

		var checkSum uint16
		dt := []byte{}

		fmt.Printf("DownloadSlotData slot: %d\n", slot)
		err = bt.DownloadSlotData(slot, 0, func(d []byte) error {
			hx, err := hex.DecodeString(string(d))
			if err != nil {
				return errors.Wrapf(err, "error decoding hex string %s", d)
			}
			dt = append(dt, hx...)
			return nil
		}, func(s string) error {
			b, err := hex.DecodeString(s)
			if err != nil {
				return errors.Wrapf(err, "could not decode crc: %s", s)
			}
			checkSum = binary.LittleEndian.Uint16(b)
			return nil
		})

		vsd := &vehdata.VehSlotData{
			DeviceMACAddress: *mac,
			Info:             vehSlotInfo,
			SlotData:         dt,
		}

		b, err := proto.Marshal(vsd)
		if err != nil {
			return errors.Wrap(err, "Marshal slot data")
		}

		return ioutil.WriteFile(fmt.Sprintf("VehSlotData%d.bin", time.Now().Unix()), b, 0644)
	}, func() error {
		return fmt.Errorf("unexpected disconnect")
	})

}
