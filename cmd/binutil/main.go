package main

import (
	"fmt"
	"os"

	"github.com/bbengfort/binutil"
	"github.com/urfave/cli/v2"
)

func main() {
	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "binutil"
	app.Version = binutil.Version()
	app.Usage = "helpers for converting to and from binary and string representations"
	app.UsageText = "binutil [-d DECODE] [-e ENCODE] [-r PATH] [INPUT]\n\n  The encoder and decoder must be one of the registered decoders;\n  to see availabe decoders:\n\nbinutil decoders\n\n  For example to convert a ulid to base64:\n\nbinutil -d ulid -e b64 01H3W3MX9A4AFNW55R0MNMQR6Y"
	app.Action = handler
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
		},
		&cli.BoolFlag{
			Name:    "binary",
			Aliases: []string{"b"},
			Usage:   "the input is binary data not a UTF-8 string",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "decoders",
			Aliases: []string{"d"},
			Usage:   "print the list of registered decoders",
			Action:  listDecoders,
		},
	}

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

	if c.String("decode") == "" || c.String("encode") == "" {
		return cli.Exit("encoder and decoder must be specified", 1)
	}

	// Handle pipeline of encoders/decoders
	var pipe *binutil.Pipeline
	if pipe, err = binutil.New(c.String("decode"), c.String("encode")); err != nil {
		return cli.Exit(err, 1)
	}

	args := c.Args()
	for i := 0; i < c.NArg(); i++ {
		var out string
		in := args.Get(i)
		if out, err = pipe.Str2Str(in); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println(out)
	}
	return nil
}

func listDecoders(c *cli.Context) error {
	names := binutil.DecoderNames()
	fmt.Println("Registered Decoders:\n====================")
	for _, name := range names {
		fmt.Printf("- %s\n", name)
	}
	return nil
}
