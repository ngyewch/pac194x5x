package main

import (
	"fmt"

	"github.com/ngyewch/pac194x5x"
	"github.com/urfave/cli/v2"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func doRead(cCtx *cli.Context) error {
	_, err := host.Init()
	if err != nil {
		return err
	}

	b, err := i2creg.Open(i2cBusFlag.Get(cCtx))
	if err != nil {
		return err
	}

	dev, err := pac194x5x.NewI2C(b, uint16(i2cAddrFlag.Get(cCtx)), voltageRatioFlag.Get(cCtx), rSenseFlag.Get(cCtx))
	if err != nil {
		return err
	}

	err = dev.RefreshV(pac194x5x.DefaultDelay)
	if err != nil {
		return err
	}

	for i := range dev.Channels() {
		fmt.Printf("[Channel #%d]\n", i+1)

		vBus, err := dev.GetVBus(i)
		if err != nil {
			return err
		}
		current, err := dev.GetCurrent(i)
		if err != nil {
			return err
		}
		power, err := dev.GetVPower(i)
		if err != nil {
			return err
		}
		vBusAvg, err := dev.GetVBusAvg(i)
		if err != nil {
			return err
		}
		currentAvg, err := dev.GetCurrentAvg(i)
		if err != nil {
			return err
		}
		energy, err := dev.GetEnergy(i)
		if err != nil {
			return err
		}

		fmt.Printf("vBus: %f V\n", vBus)
		fmt.Printf("Current: %f mA\n", current)
		fmt.Printf("Power: %f W\n", power)
		fmt.Printf("vBusAvg: %f V\n", vBusAvg)
		fmt.Printf("CurrentAvg: %f mA\n", currentAvg)
		fmt.Printf("Energy: %f µWh\n", energy)
		fmt.Println()
	}

	return nil
}
