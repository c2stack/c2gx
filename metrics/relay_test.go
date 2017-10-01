package metrics

import (
	"testing"

	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/device"
	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/nodes"
)

func TestRelayNotify(t *testing.T) {
	d := device.New(&meta.StringSource{
		Streamer: func(string) (string, error) {
			return `module x {
				revision 0;
				notification update {		
				}
			}`, nil
		},
	})
	update := make(chan bool)
	n := &nodes.Basic{
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			go func() {
				<-update
				r.Send(nil)
			}()
			return func() error {
				return nil
			}, nil
		},
	}
	if err := d.Add("x", n); err != nil {
		t.Error(err)
	}
	dm := device.NewMap()
	dm.Add("dev0", d)

	queue := make(chan Metric)
	r := &Relay{devices: dm, Sink: queue}
	o := r.Source()
	o.Path = "update"
	o.Module = "x"
	o.Device = "dev0"
	r.Script = `
	  tag("i", "1")
	  send()
	`
	err := r.SetSource(o)
	if err != nil {
		t.Fatal(err)
	}
	update <- true
	m := <-queue
	c2.AssertEqual(t, "1", m.Tags["i"])
	t.Log("got one")
}
