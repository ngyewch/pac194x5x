package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	defaultVoltageRatios = []float64{1, 1, 1, 1}
	defaultRSenses       = []float64{0.004, 0.004, 0.004, 0.004}

	i2cBusFlag = &cli.StringFlag{
		Name:     "i2c-bus",
		Usage:    "I2C bus",
		Required: true,
		EnvVars:  []string{"I2C_BUS"},
	}
	i2cAddrFlag = &cli.UintFlag{
		Name:     "i2c-addr",
		Usage:    "I2C addr",
		Required: true,
		EnvVars:  []string{"I2C_ADDR"},
	}
	voltageRatioFlag = &cli.Float64SliceFlag{
		Name:    "voltage-ratio",
		Usage:   "voltage ratio",
		Value:   cli.NewFloat64Slice(1, 1, 1, 1),
		EnvVars: []string{"VOLTAGE_RATIO"},
	}
	rSenseFlag = &cli.Float64SliceFlag{
		Name:    "rsense",
		Usage:   "RSense",
		Value:   cli.NewFloat64Slice(0.004, 0.004, 0.004, 0.004),
		EnvVars: []string{"RSENSE"},
	}

	app = &cli.App{
		Name:  "pac194x5x",
		Usage: "PAC194x/5x CLI",
		Flags: []cli.Flag{
			i2cBusFlag,
			i2cAddrFlag,
			voltageRatioFlag,
			rSenseFlag,
		},
		Commands: []*cli.Command{
			{
				Name:   "read",
				Usage:  "read",
				Action: doRead,
			},
		},
	}
)

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo != nil {
		app.Version = buildInfo.Main.Version
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
