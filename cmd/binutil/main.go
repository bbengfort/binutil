package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/bbengfort/binutil"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"
)

const (
	// Pretty print the output rather than encoding it (usually as a table).
	pretty = "pretty"
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
		{
			Name:   "ulid",
			Usage:  "generate a new ulid",
			Action: makeULID,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "encoder",
					Aliases: []string{"e"},
					Usage:   "the encoder to display the ulid repr in",
					Value:   pretty,
				},
				&cli.BoolFlag{
					Name:    "no-newline",
					Aliases: []string{"n"},
					Usage:   "omit newline, useful for use with pbcopy (ignored for encoder=pretty)",
				},
			},
		},
		{
			Name:   "uuid",
			Usage:  "generate a new uuid",
			Action: makeUUID,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "encoder",
					Aliases: []string{"e"},
					Usage:   "the encoder to display the uuid repr in",
					Value:   pretty,
				},
				&cli.BoolFlag{
					Name:    "no-newline",
					Aliases: []string{"n"},
					Usage:   "omit newline, useful for use with pbcopy (ignored for encoder=pretty)",
				},
			},
		},
		{
			Name:   "rand",
			Usage:  "generate a new random byte array",
			Action: makeRand,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "the number of bytes to generate",
					Value:   16,
				},
				&cli.StringFlag{
					Name:    "encoder",
					Aliases: []string{"e"},
					Usage:   "the encoder to display the bytes in",
					Value:   "base64",
				},
				&cli.BoolFlag{
					Name:    "no-newline",
					Aliases: []string{"n"},
					Usage:   "omit newline, useful for use with pbcopy",
				},
			},
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

func makeULID(c *cli.Context) error {
	// TODO: add timestamp parsing and other features from native ulid command!
	uu := ulid.Make()
	if encoder := c.String("encoder"); encoder != pretty {
		pipe, err := binutil.New(encoder)
		if err != nil {
			return cli.Exit(err, 1)
		}

		out, err := pipe.Bin2Str(uu.Bytes())
		if err != nil {
			return cli.Exit(err, 1)
		}

		if !c.Bool("no-newline") {
			out += "\n"
		}

		fmt.Print(out)
		return nil
	}

	// Pretty print a table with the ULID, timestamp, hex, and b64 encodings
	multi, err := binutil.NewMulti("hex", "b64")
	if err != nil {
		return cli.Exit(err, 1)
	}

	out := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.AlignRight|tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(out, "ULID\t%s\t\n", uu.String())
	fmt.Fprintf(out, "Time\t%s\t\n", ulid.Time(uu.Time()).Format(time.RFC3339Nano))
	fmt.Fprintf(out, "Hex Bytes\t%s\t\n", multi.MustBin2Str("hex", uu[:]))
	fmt.Fprintf(out, "b64 Bytes\t%s\t\n", multi.MustBin2Str("b64", uu[:]))
	out.Flush()
	return nil
}

func makeUUID(c *cli.Context) error {
	uu, err := uuid.NewRandom()
	if err != nil {
		return cli.Exit(err, 9)
	}

	if encoder := c.String("encoder"); encoder != pretty {
		pipe, err := binutil.New(encoder)
		if err != nil {
			return cli.Exit(err, 1)
		}

		out, err := pipe.Bin2Str(uu[:])
		if err != nil {
			return cli.Exit(err, 1)
		}

		if !c.Bool("no-newline") {
			out += "\n"
		}

		fmt.Print(out)
		return nil
	}

	// Pretty print a table with the UUID, hex, and b64 encodings
	multi, err := binutil.NewMulti("hex", "b64")
	if err != nil {
		return cli.Exit(err, 1)
	}

	out := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.AlignRight|tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(out, "ULID\t%s\t\n", uu.String())
	fmt.Fprintf(out, "Hex Bytes\t%s\t\n", multi.MustBin2Str("hex", uu[:]))
	fmt.Fprintf(out, "b64 Bytes\t%s\t\n", multi.MustBin2Str("b64", uu[:]))
	out.Flush()
	return nil
}

func makeRand(c *cli.Context) error {
	data := make([]byte, c.Int("size"))
	if _, err := rand.Read(data); err != nil {
		return cli.Exit(err, 9)
	}

	pipe, err := binutil.New(c.String("encoder"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	out, err := pipe.Bin2Str(data)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if !c.Bool("no-newline") {
		out += "\n"
	}

	fmt.Print(out)
	return nil
}
