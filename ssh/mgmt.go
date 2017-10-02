package ssh

import (
	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/nodes"
)

func Manage(s *Service) node.Node {
	o := s.Options()
	return &nodes.Extend{
		Base: nodes.ReflectChild(&o),
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return s.Apply(o)
		},
	}
}
