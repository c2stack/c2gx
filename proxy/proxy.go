package proxy

import (
	"bytes"
	"github.com/c2g/meta"
	"github.com/c2g/node"
	"io"
	"strings"
)

type EndpointRequest func(method string, url string, payload io.Reader) (io.ReadCloser, error)

// merges node operations between a local config node and a copy of remote operational data.
type proxy struct {
	stripPathPrefix string
	method          string
	edits           map[string]interface{}
	onCommit        func() error
	onRequest       EndpointRequest
}

func (self *proxy) proxy(config node.Node, operational node.Node) node.Node {
	self.edits = make(map[string]interface{}, 0)
	return self.container(config, operational, true, self.edits)
}

func (self *proxy) url(p *node.Path) string {
	return "restconf/" + strings.TrimPrefix(p.String(), self.stripPathPrefix)
}

func (self *proxy) container(config node.Node, remote node.Node, firstLevel bool, edits map[string]interface{}) node.Node {
	return &node.MyNode{
		Label: "proxy",
		OnSelect: func(r node.ContainerRequest) (node.Node, error) {
			if r.Target.Tail != nil {
				return r.Selection.Node(), nil
			}
			var err error
			var remoteChild, configChild node.Node
			isconfig := r.Selection.IsConfig(r.Meta)
			if config != nil {
				if configChild, err = config.Select(r); err != nil || (configChild == nil && isconfig) {
					return nil, err
				}
			}
			if remote != nil {
				readOnlyRequest := r
				readOnlyRequest.New = false
				if remoteChild, err = remote.Select(readOnlyRequest); err != nil {
					return nil, err
				}
			}
			if configChild == nil && remoteChild == nil {
				return nil, nil
			}

			if r.New && firstLevel {
				self.method = "POST"
			} else {
				self.method = "PUT"
			}
			newpos := make(map[string]interface{})
			edits[r.Meta.GetIdent()] = newpos
			return self.container(configChild, remoteChild, false, newpos), err
		},
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			if r.Target.Tail != nil {
				return r.Selection.Node(), r.Key, nil
			}
			var err error
			var key = r.Key
			var remoteChild, configChild node.Node
			isconfig := r.Selection.IsConfig(r.Meta)
			if remote != nil {
				if remoteChild, key, err = remote.Next(r); err != nil {
					return nil, nil, err
				}
			}
			if config != nil && isconfig {
				if configChild, key, err = config.Next(r); err != nil {
					return nil, nil, err
				}
			}
			if remoteChild == nil && configChild == nil {
				return nil, nil, nil
			}

			newpos := make(map[string]interface{})
			item := edits[r.Meta.GetIdent()]
			if item == nil {
				list := make([]interface{}, 1)
				list[0] = newpos
				item = list
			} else {
				list := item.([]interface{})
				list = append(list, newpos)
				item = list
			}
			edits[r.Meta.GetIdent()] = item
			return self.container(configChild, remoteChild, false, newpos), key, err
		},
		OnEvent: func(s *node.Selection, e node.Event) error {
			switch e.Type {
			case node.DELETE:
				if err := self.dispatch(config, s, e); err != nil {
					return err
				}
				if _, err := self.onRequest("DELETE", self.url(s.Path()), nil); err != nil {
					return err
				}
				return self.onCommit()
			case node.END_TREE_EDIT:
				e.StopPropagation()
				if err := self.dispatch(config, s, e); err != nil {
					return err
				}
				var buf bytes.Buffer
				js := node.NewJsonWriter(&buf).Node()
				b := node.NewBrowser2(s.Meta().(meta.MetaList), node.MapNode(self.edits))
				if err := b.Root().Selector().InsertInto(js).LastErr; err != nil {
					return err
				}
				if _, err := self.onRequest(self.method, self.url(s.Path()), &buf); err != nil {
					return err
				}

				return self.onCommit()
			default:
				if err := self.dispatch(remote, s, e); err != nil {
					return err
				}
				if err := self.dispatch(config, s, e); err != nil {
					return err
				}
			}
			return nil
		},
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) (err error) {
			if r.Write {
				edits[r.Meta.GetIdent()] = hnd.Val.Value()
				err = config.Field(r, hnd)
			} else {
				if config != nil {
					err = config.Field(r, hnd)
				}
				if remote != nil && hnd.Val == nil && err == nil {
					err = remote.Field(r, hnd)
				}
			}
			return
		},
		OnAction: func(r node.ActionRequest) (node.Node, error) {
			var buf bytes.Buffer
			js := node.NewJsonWriter(&buf).Node()
			if err := r.Input.Selector().InsertInto(js).LastErr; err != nil {
				return nil, err
			}
			body, err := self.onRequest("POST", self.url(r.Selection.Path()), &buf)
			if err != nil {
				return nil, err
			}
			var output node.Node
			if body != nil {
				output = node.NewJsonReader(body).Node()
			}
			return output, nil
		},
		OnChoose: func(sel *node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error) {
			if config != nil {
				m, err = config.Choose(sel, choice)
			}
			if remote != nil && m == nil && err == nil {
				m, err = config.Choose(sel, choice)
			}
			return
		},
	}
}

func (self *proxy) dispatch(n node.Node, s *node.Selection, e node.Event) error {
	if n != nil {
		return n.Event(s, e)
	}
	return nil
}
