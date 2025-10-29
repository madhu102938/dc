package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name: "clone",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "exclude",
						Aliases: []string{"e"},
						Usage:   "Regex pattern to exclude the files/folders",
						Value:   "a^", // Never matches anything
					},
					&cli.BoolFlag{
						Name:    "clipboard",
						Aliases: []string{"c"},
						Usage:   "Copy the script to clipboard",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:  "include-hidden",
						Usage: "include dot(.) files also",
						Value: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					if !c.Args().Present() {
						return errors.New("No arguments given for `clone`")
					}

					pattern := c.String("exclude")
					reg, err := regexp.Compile(pattern)
					if err != nil {
						fmt.Println("ERROR: invalid pattern", pattern, "", "considering empty pattern")
						reg, _ = regexp.Compile("a^") // Never matches anything
					}

					traverseDirectories(c.Args().Slice(), reg, c.Bool("clipboard"), c.Bool("include-hidden"))

					return nil
				},
				Usage: "clones directories and files given as arguments and creates a script",
			},
		},
		Usage: "A cli tool for creating a script to reproduce directories",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
