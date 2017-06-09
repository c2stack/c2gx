package metrics

import "testing"
import "github.com/c2stack/c2g/device"
import "github.com/c2stack/c2g/examples/car"
import "github.com/c2stack/c2g/meta"
import "github.com/c2stack/c2g/c2"

func Test_RelayNotify(t *testing.T) {
	d := device.New(&meta.FileStreamSource{Root: "../../c2g/examples/car"})
	c := car.New()
	c.Speed = 1
	d.Add("car", car.Node(c))
	dm := device.NewMap()
	dm.Add("dev0", d)

	queue := make(chan Metric)
	r := &Relay{devices: dm, Sink: queue}
	o := r.Source()
	o.Path = "update"
	o.Module = "car"
	o.Device = "dev0"
	r.Script = `
	  tag("i", "1")
	  send()
	`
	err := r.SetSource(o)
	if err != nil {
		t.Fatal(err)
	}
	c.Start()
	m := <-queue
	if err := c2.CheckEqual("1", m.Tags["i"]); err != nil {
		t.Error(err)
	}
	t.Log("got one")
}
