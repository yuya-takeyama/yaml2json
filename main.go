package main

import (
	"fmt"
	"io"
	"os"

	flags "github.com/jessevdk/go-flags"
	jsoniter "github.com/json-iterator/go"
	"github.com/yuya-takeyama/argf"
	"gopkg.in/yaml.v3"
)

const appName = "yaml2json"

var (
	version   = ""
	gitCommit = ""
)

type options struct {
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

var opts options

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = appName
	parser.Usage = "[OPTIONS] FILES..."

	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				parser.WriteHelp(os.Stderr)

				return
			}
		}

		errorf("flag parse error: %s", err)
		os.Exit(1)
	}

	r, err := argf.From(args)
	if err != nil {
		errorf("file loading error: %s", err)
		os.Exit(1)
	}

	err = yaml2json(r, os.Stdout, opts)
	if err != nil {
		errorf("error: %s", err)
		os.Exit(1)
	}
}

const lf = byte('\n')

func yaml2json(r io.Reader, stdout io.Writer, opts options) error {
	if opts.ShowVersion {
		_, _ = io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", appName, version, gitCommit))

		return nil
	}

	decoder := yaml.NewDecoder(r)

	var (
		json = jsoniter.ConfigCompatibleWithStandardLibrary
		d    interface{}
	)

	for {
		if err := decoder.Decode(&d); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var (
			b       []byte
			jsonErr error
		)

		b, jsonErr = json.Marshal(d)
		if jsonErr != nil {
			return jsonErr
		}

		stdout.Write(append(b, lf))
	}

	return nil
}

func errorf(message string, args ...interface{}) {
	subMessage := fmt.Sprintf(message, args...)
	_, _ = fmt.Fprintf(os.Stderr, "yaml2json: %s\n", subMessage)
}
