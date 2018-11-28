package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"incubator/config"
)

func main() {
	recover := flag.Bool("recover", false, "recover")
	flag.Parse()

	file, e := ioutil.ReadFile("./start.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	cfg := config.Config{}
	json.Unmarshal(file, &cfg)
	cfg.Layer.Id = 0
	cfg.Layer.Recover = *recover

	cfg.Process().Test()
}