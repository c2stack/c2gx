package host

import (
	"syscall"

	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/nodes"
	"github.com/c2stack/c2g/val"
	"github.com/cloudfoundry/gosigar"
)

func Manage() node.Node {
	return &nodes.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "cpu":
				return CpusNode()
			case "fs":
				return FilesystemsNode()
			case "ram":
				return MemoryNode()
			case "swap":
				return SwapNode()
			case "proc":
				return ProcsNode()
			}
			return nil, nil
		},
	}
}

func ProcsNode() (node.Node, error) {
	procs := sigar.ProcList{}
	if err := procs.Get(); err != nil {
		return nil, err
	}
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			key := r.Key
			var n node.Node
			if key != nil {
				for _, pid := range procs.List {
					if key[0].Value().(int) == pid {
						n = ProcNode(pid)
						break
					}
				}
			} else {
				row := int(r.Row)
				if row < len(procs.List) {
					pid := procs.List[row]
					n = ProcNode(pid)
					var err error
					key, err = node.NewValues(r.Meta.KeyMeta(), pid)
					if err != nil {
						return nil, nil, err
					}
				}
			}
			if n != nil {
				return n, key, nil
			}
			return nil, nil, nil
		},
	}, nil

}

func ignoreUnavail(err error) error {
	// some percentage of the process return no information on Mac, could be permission
	// related.
	if err == syscall.ENOMEM {
		return nil
	}
	return err
}

func ProcNode(pid int) node.Node {
	return &nodes.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "state":
				state := &sigar.ProcState{}
				if err := ignoreUnavail(state.Get(pid)); err != nil {
					return nil, err
				}
				return ProcStateNode(*state), nil
			case "mem":
				mem := sigar.ProcMem{}
				if err := ignoreUnavail(mem.Get(pid)); err != nil {
					return nil, err
				}
				return nodes.ReflectChild(&mem), nil
			case "time":
				time := sigar.ProcTime{}
				if err := ignoreUnavail(time.Get(pid)); err != nil {
					return nil, err
				}
				return nodes.ReflectChild(&time), nil
			}
			return nil, nil
		},
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) (err error) {
			switch r.Meta.GetIdent() {
			case "pid":
				hnd.Val = val.Int32(pid)
			}
			return nil
		},
	}
}

func ProcStateNode(state sigar.ProcState) node.Node {
	return &nodes.Extend{
		Base: nodes.ReflectChild(&state),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) (err error) {
			switch r.Meta.GetIdent() {
			case "state":
				e := r.Meta.GetDataType().EnumerationRef
				switch state.State {
				case sigar.RunStateSleep:
					hnd.Val, _ = e.ById(0)
				case sigar.RunStateRun:
					hnd.Val, _ = e.ById(1)
				case sigar.RunStateStop:
					hnd.Val, _ = e.ById(2)
				case sigar.RunStateZombie:
					hnd.Val, _ = e.ById(3)
				case sigar.RunStateIdle:
					hnd.Val, _ = e.ById(4)
				case sigar.RunStateUnknown:
					hnd.Val, _ = e.ById(5)
				}
				return nil
			}
			return p.Field(r, hnd)
		},
	}
}

func SwapNode() (node.Node, error) {
	swap := sigar.Swap{}
	if err := swap.Get(); err != nil {
		return nil, err
	}
	return nodes.ReflectChild(&swap), nil
}

func MemoryNode() (node.Node, error) {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}
	return nodes.ReflectChild(&mem), nil
}

func FilesystemsNode() (node.Node, error) {
	fs := &sigar.FileSystemList{}
	if err := fs.Get(); err != nil {
		return nil, err
	}
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			key := r.Key
			var n node.Node
			if key != nil {
				for _, fs := range fs.List {
					if key[0].String() == fs.DirName {
						n = FileSystemNode(fs)
						break
					}
				}
			} else {
				row := int(r.Row)
				if row < len(fs.List) {
					n = FileSystemNode(fs.List[row])
					key, _ = node.NewValues(r.Meta.KeyMeta(), fs.List[row].DirName)
				}
			}
			if n != nil {
				return n, key, nil
			}
			return nil, nil, nil
		},
	}, nil
}

func FileSystemNode(fss sigar.FileSystem) node.Node {
	return &nodes.Extend{
		Base: nodes.ReflectChild(&fss),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "usage":
				return FileSystemUsageNode(fss.DirName)
			}
			return p.Child(r)
		},
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) (err error) {
			switch r.Meta.GetIdent() {
			case "flags":
				hnd.Val = val.Int32(int(fss.Flags))
			default:
				err = p.Field(r, hnd)
			}
			return
		},
	}
}

func FileSystemUsageNode(dirName string) (node.Node, error) {
	var usage sigar.FileSystemUsage
	if err := usage.Get(dirName); err != nil {
		return nil, err
	}
	return nodes.ReflectChild(&usage), nil
}

func CpusNode() (node.Node, error) {
	cpus := &sigar.CpuList{}
	if err := cpus.Get(); err != nil {
		return nil, err
	}
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			key := r.Key
			var n node.Node
			if key != nil {
				if key[0].Value().(int) < len(cpus.List) {
					n = CpuNode(key[0].Value().(int), cpus.List[key[0].Value().(int)])
				}
			} else {
				id := int(r.Row)
				if id < len(cpus.List) {
					n = CpuNode(id, cpus.List[id])
					key, _ = node.NewValues(r.Meta.KeyMeta(), id)
				}
			}
			if n != nil {
				return n, key, nil
			}
			return nil, nil, nil
		},
	}, nil
}

func CpuNode(id int, cpu sigar.Cpu) node.Node {
	return &nodes.Extend{
		Base: nodes.ReflectChild(&cpu),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) (err error) {
			switch r.Meta.GetIdent() {
			case "id":
				hnd.Val = val.Int32(id)
			default:
				err = p.Field(r, hnd)
			}
			return
		},
	}
}
