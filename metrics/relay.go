package metrics

import (
	"time"

	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/conf"
	"github.com/c2stack/c2g/node"
	"github.com/mattn/anko/vm"
)

type Relay struct {
	Name     string
	Database string
	Script   string
	devices  conf.ServiceLocator
	src      NotifySource
	tags     map[string]string
	closer   node.NotifyCloser
	Sink     chan<- Metric
}

type NotifySource struct {
	Device string
	Module string
	Path   string
}

func (self *Relay) Source() NotifySource {
	return self.src
}

func (self *Relay) SetSource(src NotifySource) error {
	if self.src == src {
		return nil
	}
	self.src = src
	d, err := self.devices.Device(src.Device)
	if err != nil {
		return err
	}
	b, err := d.Browser(src.Module)
	if err != nil {
		return err
	}
	if b == nil {
		return c2.NewErrC("No browser found with module : "+src.Module, 404)
	}
	s := b.Root().Find(src.Path)
	if s.LastErr != nil {
		return s.LastErr
	}
	self.Close()
	self.closer, err = s.Notifications(self.handleEvent)
	return err
}

func (self *Relay) NewMetric() *Metric {
	p := &Metric{
		Name:     self.Name,
		Database: self.Database,
		Tags:     make(map[string]string),
		Fields:   make(map[string]interface{}),
		Time:     time.Now(),
	}
	for tag, value := range self.tags {
		p.Tags[tag] = value
		p.Tags["device"] = self.src.Device
	}
	return p
}

func (self *Relay) Close() {
	if self.closer != nil {
		self.closer()
		self.closer = nil
	}
}

type Metric struct {
	Database string
	Name     string
	Tags     map[string]string
	Fields   map[string]interface{}
	Time     time.Time
}

type NodeBuiltins struct{}

func LoadNodeBuiltins(env *vm.Env) {
	x := NodeBuiltins{}
	env.Define("get", x.Get)
	env.Define("find", x.Find)
}

func (NodeBuiltins) Get(s node.Selection, ident string) interface{} {
	v, err := s.GetValue(ident)
	if err != nil {
		panic(err)
	}
	return v.Value()
}

func (NodeBuiltins) Find(s node.Selection, path string) node.Selection {
	n := s.Find(path)
	if n.LastErr != nil {
		panic(n.LastErr)
	}
	return n
}

func (self *Relay) handleEvent(event node.Selection) {
	env := vm.NewEnv()
	p := self.NewMetric()
	env.Define("event", event)
	env.Define("tag", func(tag string, value string) {
		p.Tags[tag] = value
	})
	env.Define("value", func(field string, value interface{}) {
		p.Fields[field] = value
	})
	LoadNodeBuiltins(env)
	//core.LoadAllBuiltins(env)
	env.Define("send", func() {
		self.Sink <- *p
		p = self.NewMetric()
	})
	_, err := env.Execute(self.Script)
	if err != nil {
		panic(err)
	}
}
