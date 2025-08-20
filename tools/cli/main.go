package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
)

var (
	defaultVoltageRatios = []float64{1, 1, 1, 1}
	defaultRSenses       = []float64{0.004, 0.004, 0.004, 0.004}

	i2cBusFlag = &cli.StringFlag{
		Name:     "i2c-bus",
		Usage:    "I2C bus",
		Required: true,
		Sources:  cli.EnvVars("I2C_BUS"),
	}
	i2cAddrFlag = &cli.UintFlag{
		Name:     "i2c-addr",
		Usage:    "I2C addr",
		Required: true,
		Sources:  cli.EnvVars("I2C_ADDR"),
	}
	voltageRatioFlag = &cli.Float64SliceFlag{
		Name:    "voltage-ratio",
		Usage:   "voltage ratio",
		Value:   []float64{1, 1, 1, 1},
		Sources: cli.EnvVars("VOLTAGE_RATIO"),
	}
	rSenseFlag = &cli.Float64SliceFlag{
		Name:    "rsense",
		Usage:   "RSense",
		Value:   []float64{0.004, 0.004, 0.004, 0.004},
		Sources: cli.EnvVars("RSENSE"),
	}

	app = &cli.Command{
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

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
