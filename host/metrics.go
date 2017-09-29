package host

import (
	"time"

	"github.com/cloudfoundry/gosigar"
)

type Metrics struct {
	SampleRate int
	Sigar      sigar.Sigar
	Cpu        sigar.Cpu
	closer     func()
}

func (self *Metrics) Start() {
	self.Sigar = &sigar.ConcreteSigar{}
	cpus, stop := self.Sigar.CollectCpuStats(time.Duration(self.SampleRate) * time.Millisecond)
	self.closer = func() {
		stop <- struct{}{}
	}
	for {
		if cpu, open := <-cpus; !open {
			return
		} else {
			self.Cpu = cpu
		}
	}
}
