package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"incubator/config"
	"incubator/message"
	"net"
	"incubator/serialization"
)

type node struct {
	Address string `json:"Address"`
}

func main() {
	layer := flag.Int("layer", 0, "layer id")
	nodesFile := flag.String("nodes-file", "", "nodes file")
	groundLayerFile := flag.String("ground-layer-file", "", "ground layer file")
	layerFile := flag.String("layer-file", "", "layer file")
	recover := flag.Bool("recover", false, "recover")
	flag.Parse()

	if *layer  < 1 {
		fmt.Printf("illegal layer(<1): %d", *layer)
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

	groundLayer := groudCfg.GetLayer()

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
	cfg.Layer.Id = int32(*layer)
	cfg.Layer.Recover = *recover

	file, e = ioutil.ReadFile(*nodesFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	nodesCfg := make([]node, 0, 0)
	json.Unmarshal(file, &nodesCfg)

	for _, node := range nodesCfg {
		go func(){
			conn, e := net.Dial("tcp", node.Address)
			if e != nil {
				fmt.Printf("Connection error: %v\n", e)
				os.Exit(1)
			}

			msg := groundLayer.GetMessageFromClass(message.PullUpMessageClassName).
				Right().Replicate().
				Right().(*message.PullUpMessage)

			msg.ToServer()
			msg.SetAddr(node.Address)
			msg.SetCfg(&cfg)

			_, e = conn.Write(serialization.Marshal(msg))
			if e != nil {
				fmt.Printf("Send message error: %v\n", e)
				os.Exit(1)
			}

			conn.r
		}()
	}
}
