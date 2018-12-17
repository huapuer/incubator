package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/incubator/config"
	_ "github.com/incubator/host"
	"github.com/incubator/interfaces"
	_ "github.com/incubator/layer"
	"github.com/incubator/message"
	_ "github.com/incubator/topo"
	"io/ioutil"
	"net"
	"os"
	"sync"
)

type node struct {
	Address string `json:"Address"`
}

func main() {
	nodesFile := flag.String("nodes-file", "./nodes.json", "nodes file")
	groundLayerFile := flag.String("ground-layer-file", "./start.json", "ground layer file")
	layerFile := flag.String("layer-file", "./pullup.json", "layer file")
	starMode := flag.String("mode", "new", "start mode")

	flag.Parse()

	file, e := ioutil.ReadFile(*groundLayerFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	groundCfg := config.Config{}
	e = json.Unmarshal(file, &groundCfg)
	if e != nil {
		fmt.Printf("Config error: %v\n", e)
		os.Exit(1)
	}
	groundCfg.Process().Test()

	groundLayer := interfaces.GetLayer(groundCfg.Layer.Id).Right()

	file, e = ioutil.ReadFile(*layerFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	cfg := config.Config{}
	e = json.Unmarshal(file, &cfg)
	if e != nil {
		fmt.Printf("Config error: %v\n", e)
		os.Exit(1)
	}

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

	file, e = ioutil.ReadFile(*nodesFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	nodesCfg := make([]node, 0, 0)
	json.Unmarshal(file, &nodesCfg)

	wg := sync.WaitGroup{}
	for _, node := range nodesCfg {
		wg.Add(1)
		go func() {
			conn, e := net.Dial("tcp", node.Address)
			if e != nil {
				fmt.Printf("Connection error: %v\n", e)
				os.Exit(1)
			}

			msg := &message.PullUpMessage{
				Cfg: &cfg,
			}
			msg.SetLayer(int8(groundCfg.Layer.Id))
			msg.SetType(int8(groundLayer.GetMessageType(msg).Right()))

			service := groundLayer.GetService(0).Right()
			_, e = conn.Write(service.GetProtocal().Pack(msg))
			if e != nil {
				fmt.Printf("Send message error: %v\n", e)
				os.Exit(1)
			}

			service.HandleConnection(context.Background(), conn)

			conn.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}
