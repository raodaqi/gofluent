package main

import (
	"flag"
)

type Context struct {
	tag    string
	record Record
}

type Record struct {
	timestamp int64
	data      map[string]string
}

type Configure struct {
	Inputs_config  []interface{}
	Outputs_config []interface{}
}

var config Configure

func init() {
	c := flag.String("c", "fluent.conf", "config filepath")
	flag.Parse()

	configure, _ := ParseConfig(nil, *c)

	for _, v := range configure.Root.Elems {
		if v.Name == "source" {
			config.Inputs_config = append(config.Inputs_config, v.Attrs)
		} else if v.Name == "match" {
			config.Outputs_config = append(config.Outputs_config, v.Attrs)
		}
	}
}

func main() {

	ctxInput := make(chan Context, 10)
	ctxOutput := make(chan Context, 10)

	go NewInputs(ctxInput)
	go NewOutputs(ctxOutput)

	for {
		ctx := <-ctxInput
		ctxOutput <- ctx
	}
}
