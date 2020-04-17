package main

import (
	"fmt"
	"os"

	"github.com/farzadmf/termask/pkg/mask"
	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "mode",
			Usage:    "(tf|json) mode determines the type of the input",
			Aliases:  []string{"m"},
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "property",
			Usage:   "property whose value we want to mask (can be specified multiple times)",
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "ignore-case",
			Usage:   "case insensitive match",
			Aliases: []string{"i"},
		},
	}

	app = cli.App{
		Name:  "termask",
		Usage: "Mask values in the terminal",
		Flags: flags,
		Action: func(c *cli.Context) error {
			mode := c.String("m")
			ignoreCase := c.Bool("i")
			properties := c.StringSlice("p")

			config := mask.NewConfig(os.Stdin, os.Stdout)

			switch mode {
			case "tf":
				masker := mask.NewTFMasker(properties, ignoreCase)
				masker.Mask(config)

				return nil
			case "json":
				masker := mask.NewJSONMasker(properties, ignoreCase)
				masker.Mask(config)

				return nil

			default:
				return cli.NewExitError(fmt.Sprintf("unknown mode: %q", mode), 1)
			}
		},
	}
)
