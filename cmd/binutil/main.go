package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/bbengfort/binutil"
	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "binutil"
	app.Version = binutil.Version()
	app.Usage = "helpers for converting to and from binary and string representations"
	app.UsageText = "binutil [-d DECODE] [-e ENCODE] [-r PATH] [INPUT]"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "read data from the specified path on disk",
		},
		&cli.StringFlag{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "the format to decode the input from",
		},
		&cli.StringFlag{
			Name:    "encode",
			Aliases: []string{"e"},
			Usage:   "the format to encode the input to",
			Value:   "string",
		},
		&cli.BoolFlag{
			Name:    "binary",
			Aliases: []string{"b"},
			Usage:   "the input is binary data not a UTF-8 string",
		},
	}
	app.Action = handler

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(2)
	}
}

func handler(c *cli.Context) (err error) {
	if c.NArg() > 0 && c.String("read") != "" {
		return cli.Exit("cannot specify input arguments and a path to read from", 1)
	}

	// TODO: handle reading from a file
	if c.String("read") != "" {
		return cli.Exit("reading from a file not implemented yet", 3)
	}

	// TODO: handle reading from stdin
	if c.NArg() == 0 {
		return cli.Exit("reading from stdin not implemented yet", 3)
	}

	// TODO: handle binary input
	if c.Bool("binary") {
		return cli.Exit("reading binary input not implemented yet", 3)
	}

	args := c.Args()
	for i := 0; i < c.NArg(); i++ {
		in := args.Get(i)

		var uu ulid.ULID
		if uu, err = ulid.Parse(in); err != nil {
			return cli.Exit(fmt.Errorf("could not parse %q: %w", in, err), 1)
		}

		fmt.Println(base64.StdEncoding.EncodeToString(uu.Bytes()))
	}
	return nil
}
