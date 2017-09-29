package metrics

import (
	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/nodes"
	"github.com/c2stack/c2g/val"
)

func InfluxNode(influx *InfluxSink) node.Node {
	o := influx.Options()
	return &nodes.Extend{
		Base: nodes.ReflectChild(&o),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "relay":
				return relayListNode(influx, node.NewIndex(influx.relays)), nil
			}
			return nil, nil
		},
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return influx.ApplyOptions(o)
		},
	}
}

func relayListNode(mgr Manager, relayIndex *node.Index) node.Node {
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var relay *Relay
			key := r.Key
			if r.New {
				name := r.Key[0].String()
				relay := &Relay{
					Name: name,
				}
				mgr.AddRelay(relay)
			}
			if key != nil {
				name := r.Key[0].String()
				relay = mgr.GetRelay(name)
				if relay != nil && r.Delete {
					relay.Close()
					return nil, nil, nil
				}
			} else {
				if v := relayIndex.NextKey(r.Row); v != node.NO_VALUE {
					name := v.String()
					if relay = mgr.GetRelay(name); relay != nil {
						var err error
						if key, err = node.NewValues(r.Meta.KeyMeta(), name); err != nil {
							return nil, nil, err
						}
					}
				}
			}
			if relay != nil {
				return relayNode(relay), key, nil
			}
			return nil, nil, nil
		},
	}
}

func relayNode(relay *Relay) node.Node {
	return &nodes.Extend{
		Base: nodes.ReflectChild(relay),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "source":
				return relaySourceNode(relay), nil
			case "tag":
				if r.Delete {
					relay.tags = nil
				} else if r.New {
					relay.tags = make(map[string]string)
				}
				if relay.tags != nil {
					return tagListNode(relay.tags), nil
				}
			default:
				return p.Child(r)
			}
			return nil, nil
		},
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.GetIdent() {
			case "update":

			}
			return nil, nil
		},
	}
}

func tagListNode(tags map[string]string) node.Node {
	index := node.NewIndex(tags)
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var tag string
			var value string
			key := r.Key
			if r.New {
				tag = r.Key[0].String()
			}
			if key != nil {
				tag = key[0].String()
				if r.Delete {
					delete(tags, tag)
				}
				value = tags[tag]
			} else {
				if v := index.NextKey(r.Row); v != node.NO_VALUE {
					tag = v.String()
					if value = tags[tag]; value != "" {
						var err error
						if key, err = node.NewValues(r.Meta.KeyMeta(), value); err != nil {
							return nil, nil, err
						}
					}
				}
			}
			if value != "" {
				return tagNode(tag, tags), key, nil
			}
			return nil, nil, nil
		},
	}
}

func tagNode(tag string, tags map[string]string) node.Node {
	return &nodes.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "name":
				if r.Write {
					v := tags[tag]
					delete(tags, tag)
					tags[hnd.Val.String()] = v
				} else {
					hnd.Val = val.String(tag)
				}
			case "value":
				if r.Write {
					tags[tag] = hnd.Val.String()
				} else {
					hnd.Val = val.String(tags[tag])
				}
			}
			return nil
		},
	}
}

func fieldListNode(fields map[string]interface{}) node.Node {
	index := node.NewIndex(fields)
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var field string
			var value interface{}
			key := r.Key
			if r.New {
				field = r.Key[0].String()
			}
			if key != nil {
				field = key[0].String()
				if r.Delete {
					delete(fields, field)
				}
				value = fields[field]
			} else {
				if v := index.NextKey(r.Row); v != node.NO_VALUE {
					field = v.String()
					if value = fields[field]; value != "" {
						var err error
						if key, err = node.NewValues(r.Meta.KeyMeta(), value); err != nil {
							return nil, nil, err
						}
					}
				}
			}
			if value != "" {
				return fieldNode(field, fields), key, nil
			}
			return nil, nil, nil
		},
	}
}

func fieldNode(field string, fields map[string]interface{}) node.Node {
	return &nodes.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "name":
				if r.Write {
					v := fields[field]
					delete(fields, field)
					fields[hnd.Val.String()] = v
				} else {
					hnd.Val = val.String(field)
				}
			case "value":
				if r.Write {
					fields[field] = hnd.Val.Value()
				} else {
					hnd.Val = val.Any{Thing: fields[field]}
				}
			}
			return nil
		},
	}
}

func relaySourceNode(relay *Relay) node.Node {
	src := relay.Source()
	return &nodes.Extend{
		Base: nodes.ReflectChild(&src),
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return relay.SetSource(src)
		},
	}
}
