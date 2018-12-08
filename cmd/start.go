package main

import (
	"../config"
	"../layer"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	//TODO: init log module

	starMode := flag.String("mode", "new", "start mode")
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

	switch *starMode {
	case "new":
		cfg.Layer.StartMode = config.LAYER_START_MODE_NEW
	case "recover":
		cfg.Layer.StartMode = config.LAYER_START_MODE_RECOVER
	case "reboot":
		cfg.Layer.StartMode = config.LAYER_START_MODE_REBOOT
	default:
		fmt.Printf("unknow start mode: %d\n", *starMode)
		os.Exit(1)
	}

	cfg.Process().Test()
	layer.GetLayer(cfg.Layer.Id).Right().Start()
}
