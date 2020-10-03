package main

import (
	"fmt"
	"io"
	"os"

	flags "github.com/jessevdk/go-flags"
	jsoniter "github.com/json-iterator/go"
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

	var filenames []string

	if len(args) == 0 {
		filenames = append(filenames, "-")
	} else {
		filenames = args
	}

	err = yaml2json(filenames, os.Stdout, opts)
	if err != nil {
		errorf("error: %s", err)
		os.Exit(1)
	}
}

const lf = byte('\n')

func yaml2json(filenames []string, stdout io.Writer, opts options) error {
	if opts.ShowVersion {
		_, _ = io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", appName, version, gitCommit))

		return nil
	}

	for _, filename := range filenames {
		err := handleFile(filename, stdout)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleFile(filename string, stdout io.Writer) error {
	var r io.Reader

	if filename == "-" {
		r = os.Stdin
	} else {
		var file *os.File
		var openErr error

		file, openErr = os.Open(filename)
		if openErr != nil {
			return fmt.Errorf("file loading error: %w", openErr)
		}

		defer file.Close()

		r = file
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
