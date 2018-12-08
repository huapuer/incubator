package main

import (
	"../config"
	"../layer"
	"../message"
	"../serialization"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
)

type node struct {
	Address string `json:"Address"`
}

func main() {
	layerId := flag.Int("layer", 0, "layer id")
	nodesFile := flag.String("nodes-file", "", "nodes file")
	groundLayerFile := flag.String("ground-layer-file", "", "ground layer file")
	layerFile := flag.String("layer-file", "", "layer file")
	flag.Parse()

	if *layerId < 1 {
		fmt.Printf("illegal layer(<1): %d", *layerId)
		os.Exit(1)
	}

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

	groundLayer := layer.GetLayer(groundCfg.Layer.Id).Right()

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
	cfg.Layer.Id = int32(*layerId)
	cfg.Layer.Recover = false

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

			_, e = conn.Write(serialization.Marshal(msg))
			if e != nil {
				fmt.Printf("Send message error: %v\n", e)
				os.Exit(1)
			}

			groundLayer.GetServer().HandleConnection(context.Background(), conn)

			conn.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}
