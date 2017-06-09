package metrics

import (
	"github.com/c2stack/c2g/node"
)

func InfluxNode(influx *InfluxSink) node.Node {
	o := influx.Options()
	return &node.Extend{
		Node: node.ReflectNode(&o),
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
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var relay *Relay
			key := r.Key
			if r.New {
				name := r.Key[0].Str
				relay := &Relay{
					Name: name,
				}
				mgr.AddRelay(relay)
			}
			if key != nil {
				name := r.Key[0].Str
				relay = mgr.GetRelay(name)
				if relay != nil && r.Delete {
					relay.Close()
					return nil, nil, nil
				}
			} else {
				if v := relayIndex.NextKey(r.Row); v != node.NO_VALUE {
					name := v.String()
					if relay = mgr.GetRelay(name); relay != nil {
						key = node.SetValues(r.Meta.KeyMeta(), name)
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
	return &node.Extend{
		Node: node.ReflectNode(relay),
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
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var tag string
			var value string
			key := r.Key
			if r.New {
				tag = r.Key[0].Str
			}
			if key != nil {
				tag = key[0].Str
				if r.Delete {
					delete(tags, tag)
				}
				value = tags[tag]
			} else {
				if v := index.NextKey(r.Row); v != node.NO_VALUE {
					tag = v.String()
					if value = tags[tag]; value != "" {
						key = node.SetValues(r.Meta.KeyMeta(), value)
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
	return &node.MyNode{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "name":
				if r.Write {
					v := tags[tag]
					delete(tags, tag)
					tags[hnd.Val.Str] = v
				} else {
					hnd.Val = &node.Value{Str: tag}
				}
			case "value":
				if r.Write {
					tags[tag] = hnd.Val.Str
				} else {
					hnd.Val = &node.Value{Str: tags[tag]}
				}
			}
			return nil
		},
	}
}

func fieldListNode(fields map[string]interface{}) node.Node {
	index := node.NewIndex(fields)
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var field string
			var value interface{}
			key := r.Key
			if r.New {
				field = r.Key[0].Str
			}
			if key != nil {
				field = key[0].Str
				if r.Delete {
					delete(fields, field)
				}
				value = fields[field]
			} else {
				if v := index.NextKey(r.Row); v != node.NO_VALUE {
					field = v.String()
					if value = fields[field]; value != "" {
						key = node.SetValues(r.Meta.KeyMeta(), value)
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
	return &node.MyNode{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "name":
				if r.Write {
					v := fields[field]
					delete(fields, field)
					fields[hnd.Val.Str] = v
				} else {
					hnd.Val = &node.Value{Str: field}
				}
			case "value":
				if r.Write {
					fields[field] = hnd.Val.AnyData
				} else {
					hnd.Val = &node.Value{AnyData: fields[field]}
				}
			}
			return nil
		},
	}
}

func relaySourceNode(relay *Relay) node.Node {
	src := relay.Source()
	return &node.Extend{
		Node: node.ReflectNode(&src),
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return relay.SetSource(src)
		},
	}
}
