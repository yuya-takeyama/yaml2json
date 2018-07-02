package main

import (
	"fmt"
	"io"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/json-iterator/go"
	"github.com/yuya-takeyama/argf"
	"gopkg.in/yaml.v2"
)

const AppName = "yaml2json"

type Options struct {
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = AppName
	parser.Usage = "[OPTIONS] FILES..."

	args, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	r, err := argf.From(args)
	if err != nil {
		panic(err)
	}

	err = yaml2json(r, os.Stdout, os.Stderr, opts)
	if err != nil {
		panic(err)
	}
}

const lf = byte('\n')

func yaml2json(r io.Reader, stdout io.Writer, stderr io.Writer, opts Options) error {
	if opts.ShowVersion {
		io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return nil
	}

	decoder := yaml.NewDecoder(r)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var d interface{}

	for {
		if err := decoder.Decode(&d); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("YAML decoding error: %s", err)
		}

		if b, err := json.Marshal(d); err != nil {
			log.Fatalf("JSON encoding error: %s", err)
		} else {
			stdout.Write(append(b, lf))
		}
	}

	return nil
}
