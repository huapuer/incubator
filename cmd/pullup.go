package main

import (
	"../config"
	"../message"
	"../serialization"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"incubator/layer"
	"context"
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
	recover := flag.Bool("recover", false, "recover")
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

	groudCfg := config.Config{}
	e = json.Unmarshal(file, &groudCfg)
	if e != nil {
		fmt.Printf("Config error: %v\n", e)
		os.Exit(1)
	}
	groudCfg.Process().Test()

	groundLayer := layer.GetLayer(groudCfg.Layer.Id).Right()

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
	cfg.Layer.Recover = *recover

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

			msg := &message.PullUpMessage{}
			msg.SetLayer(int8(groudCfg.Layer.Id))
			msg.SetType(int8(groundLayer.GetMessageType(msg).Right()))

			msg.SetAddr(node.Address)
			msg.SetCfg(&cfg)

			_, e = conn.Write(serialization.Marshal(msg))
			if e != nil {
				fmt.Printf("Send message error: %v\n", e)
				os.Exit(1)
			}

			groundLayer.GetServer().HandleConnection(context.Background(), conn)

			wg.Done()
		}()
	}

	wg.Wait()
}
