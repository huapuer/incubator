package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	_ "github.com/incubator/layer"
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
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Printf("unmarshal cfg failed: %v\n", err)
		os.Exit(1)
	}

	cfg.Layer.Id = 0

	fmt.Printf("%+v", cfg)

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
	interfaces.GetLayer(cfg.Layer.Id).Right().Start()
}
