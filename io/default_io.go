package io

import (
	"../common/maybe"
	"../config"
	"../layer"
	"../message"
	"../network"
	"fmt"
)

const (
	defaultIOClassName = "io.defaultIO"
)

func init() {
	RegisterIOPrototype(defaultIOClassName, &defaultIO{}).Test()
}

type defaultIO struct {
	commonIO

	inputJoints  []joint
	outputJoints []joint
}

func (this defaultIO) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeIO{}

	inputJointsCfg := config.GetAttrMapEfaceArray(attrs, "InputJoints").Right().([]map[string]interface{})
	outputJointsCfg := config.GetAttrMapEfaceArray(attrs, "OutputJoints").Right().([]map[string]interface{})

	topoSchema := config.GetAttrInt32(cfg.Layer.Attributes, "TopoSchema", config.CheckInt32GT0).Right()
	topoCfg, ok := cfg.Topos[topoSchema]
	if !ok {
		ret.Error(fmt.Errorf("topo cfg not found: %d", topoSchema))
		return ret
	}
	totalHostNum := config.GetAttrInt64(topoCfg.Attributes, "TotalHostNum", config.CheckInt64GT0).Right()

	value := &defaultIO{
		inputJoints:  make([]joint, 0, 0),
		outputJoints: make([]joint, 0, 0),
	}

	for _, jointCfg := range inputJointsCfg {
		begin := config.GetAttrInt64(jointCfg, "Begin", config.CheckInt64GT0).Right()
		end := config.GetAttrInt64(jointCfg, "End", config.CheckInt64GT0).Right()

		if begin > totalHostNum {
			ret.Error(fmt.Errorf("input begin host id exceeds host range: %d>%d", begin, totalHostNum))
			return ret
		}
		if end > totalHostNum {
			ret.Error(fmt.Errorf("input end host id exceeds host range: %d>%d", end, totalHostNum))
			return ret
		}

		j := joint{
			begin: begin,
			end:   end,
		}
		value.inputJoints = append(value.inputJoints, j)
	}

	for _, jointCfg := range outputJointsCfg {
		begin := config.GetAttrInt64(jointCfg, "Begin", config.CheckInt64GT0).Right()
		end := config.GetAttrInt64(jointCfg, "End", config.CheckInt64GT0).Right()

		if begin > totalHostNum {
			ret.Error(fmt.Errorf("output begin host id exceeds host range: %d>%d", begin, totalHostNum))
			return ret
		}
		if end > totalHostNum {
			ret.Error(fmt.Errorf("output end host id exceeds host range: %d>%d", end, totalHostNum))
			return ret
		}

		ip := config.GetAttrString(jointCfg, "IP", config.CheckStringNotEmpty).Right()
		port := config.GetAttrInt(jointCfg, "Port", config.CheckIntGT0).Right()

		clientSchema := config.GetAttrInt32(jointCfg, "ClientSchema", nil).Right()
		clientCfg, ok := cfg.Clients[clientSchema]
		if !ok {
			ret.Error(fmt.Errorf("client cfg not found: %d", clientSchema))
			return ret
		}

		c := network.GetClientPrototype(clientCfg.Class).
			Right().New(clientCfg.Attributes, cfg).(network.MaybeClient).Right()
		c.Connect(fmt.Sprintf("%s:%d", ip, port))

		j := joint{
			begin:  begin,
			end:    end,
			client: c,
		}
		value.outputJoints = append(value.outputJoints, j)
	}

	ret.Value(value)
	return ret
}

func (this defaultIO) Input(host int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	for _, joint := range this.inputJoints {
		if host >= joint.begin && host <= joint.end {
			topo := layer.GetLayer(this.layerId).Right().GetTopo()
			topo.SendToHost(host, msg).Test()

			err.Error(nil)
			return
		}
	}
	err.Error(fmt.Errorf("host not in input range: %d", host))
	return
}

func (this defaultIO) Output(host int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	for _, joint := range this.outputJoints {
		if host >= joint.begin && host <= joint.end {
			joint.client.Send(msg).Test()
			err.Error(nil)
			return
		}
	}
	err.Error(fmt.Errorf("host not in output range: %d", host))
	return
}
