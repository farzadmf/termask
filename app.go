package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "property",
			Usage:   "property to mask (can be specified multiple times)",
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "ignore-case",
			Usage:   "case insensitive match",
			Aliases: []string{"i"},
		},
	}

	app = cli.App{
		Name:  "tfmask",
		Usage: "Mask Terraform property values",
		Flags: flags,
		Action: func(c *cli.Context) error {
			ignoreCase := c.Bool("i")
			properties := c.StringSlice("p")

			matcher := NewMatcher()
			masker := NewMasker(matcher, properties, ignoreCase)
			masker.Mask(MaskConfig{
				reader: os.Stdin,
				writer: os.Stdout,
			})

			return nil
		},
	}
)
