package proxy

import (
	"github.com/c2g/node"
	"github.com/c2g/meta"
)

type Api struct {
}

func (self Api) Registrar(registrar *Registrar) node.Node {
	return &node.Extend{
		Node: node.MarshalContainer(registrar),
		OnSelect: func(p node.Node, r node.ContainerRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "endpoint":
				if self.Registrar != nil {
					return self.Endpoints(registrar), nil
				}
				return nil, nil
			}
			return p.Select(r)
		},
		OnAction: func(p node.Node, r node.ActionRequest) (output node.Node, err error) {
			switch r.Meta.GetIdent() {
			case "register":
				return self.RegisterEndpoint(registrar, r.Selection, r.Meta, r.Input)
			}
			return p.Action(r)
		},
	}
}

func (self Api) Endpoint(endpoint *Endpoint) node.Node {
	return &node.Extend{
		Node: node.MarshalContainer(endpoint),
		OnSelect: func(p node.Node, r node.ContainerRequest) (node.Node, error) {
			if r.Meta == endpoint.Meta {
				return endpoint.handleRequest(r.Target)
			}
			return p.Select(r)
		},
		OnChoose: func(p node.Node, sel *node.Selection, choice *meta.Choice) (*meta.ChoiceCase, error) {
			return choice.GetCase(endpoint.Meta.GetIdent()), nil
		},
	}
}

func (self Api) Endpoints(registrar *Registrar) node.Node {
	n := &node.MarshalMap{
		Map: registrar.Endpoints,
		OnSelectItem: func(item interface{}) node.Node {
			return self.Endpoint(item.(*Endpoint))
		},
	}
	return n.Node()
}

func (self Api) RegisterEndpoint(registrar *Registrar, sel *node.Selection, rpc *meta.Rpc, input *node.Selection) (output node.Node, err error) {
	var reg Endpoint
	regNode := node.MarshalContainer(&reg)
	if err = input.Selector().UpsertInto(regNode).LastErr; err != nil {
		return nil, err
	}
	if err = registrar.RegisterEndpoint(&reg); err != nil {
		return nil, err
	}
	return regNode, nil
}

