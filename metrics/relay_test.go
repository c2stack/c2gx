package metrics

import "testing"
import "github.com/c2stack/c2g/device"
import "github.com/c2stack/examples/car"
import "github.com/c2stack/c2g/meta"
import "github.com/c2stack/c2g/c2"

func Test_RelayNotify(t *testing.T) {
	d := device.New(&meta.FileStreamSource{Root: "../../examples/car"})
	c := car.New()
	c.Speed = 1
	d.Add("car", car.Manage(c))
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
	c2.AssertEqual(t, "1", m.Tags["i"])
	t.Log("got one")
}
