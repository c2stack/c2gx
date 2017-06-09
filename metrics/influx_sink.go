package metrics

import (
	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/device"
	client "github.com/influxdata/influxdb/client/v2"
)

type InfluxSink struct {
	Devices device.ServiceLocator
	options InfluxOptions
	conn    client.Client
	queue   chan Metric
	relays  map[string]*Relay
}

func NewInfluxSink(devices device.ServiceLocator) *InfluxSink {
	sink := &InfluxSink{
		Devices: devices,
		queue:   make(chan Metric),
		relays:  make(map[string]*Relay),
	}
	go sink.Start()
	return sink
}

func (self *InfluxSink) Options() InfluxOptions {
	return self.options
}

func (self *InfluxSink) ApplyOptions(options InfluxOptions) error {
	if self.options == options {
		return nil
	}
	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: options.Addr,
	})
	if err != nil {
		return err
	}
	self.options = options
	self.conn = conn
	return nil
}

func (self *InfluxSink) Start() {
	for {
		m := <-self.queue
		p, err := client.NewPoint(m.Name, m.Tags, m.Fields, m.Time)
		if err != nil {
			panic(err)
		}
		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  m.Database,
			Precision: "s",
		})
		bp.AddPoint(p)
		c2.Debug.Printf("point %v", p)
		if self.conn != nil {
			self.conn.Write(bp)
		}
	}
}

type InfluxOptions struct {
	Addr string
}

func (self *InfluxSink) GetRelay(name string) *Relay {
	return self.relays[name]
}

func (self *InfluxSink) RemoveRelay(name string) {
	if relay, found := self.relays[name]; found {
		relay.Close()
		delete(self.relays, name)
	}
}

func (self *InfluxSink) AddRelay(relay *Relay) {
	relay.devices = self.Devices
	relay.Sink = self.queue
	self.relays[relay.Name] = relay
}
