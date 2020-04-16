package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/farzadmf/termask/pkg/mask"
	"github.com/farzadmf/termask/pkg/match"
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

			switch mode {
			case "tf":
				m := match.NewTFMatcher()
				masker := mask.NewMasker(m, properties, ignoreCase)
				masker.Mask(mask.NewConfig(os.Stdin, os.Stdout))

				return nil
			case "json":
				m := match.NewJSONMatcher()
				names, matches := m.Hello(`"prop": "value"`)
				for i := 0; i < len(names); i++ {
					fmt.Println("name", names[i])
					if names[i] == "value" {
						matches[i] = "***"
					}
				}

				fmt.Println(strings.Join(matches[1:], ""))

				return nil

			default:
				return cli.NewExitError(fmt.Sprintf("unknown mode: %q", mode), 1)
			}
		},
	}
)
