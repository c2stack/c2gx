package registrar

import (
	"github.com/c2g/node"
	"net/http"
	"github.com/c2g/meta"
	"github.com/c2g/c2"
)

type Endpoint struct {
	Id              string
	Module          string
	module          *meta.Module
	EndpointAddress string
	Store           node.Store
}

func (self *Endpoint) Schema() (*meta.Module, error) {
	if self.module != nil {
		return self.module, nil
	}
	url := self.EndpointAddress + "meta/?noexpand"
	c2.Info.Printf("Downloading meta %s", url)
	in, err := self.DownloadJson(url)
	if err != nil {
		return nil, err
	}
	m := &meta.Module{}
	if err = node.SelectModule(m, false).Root().Selector().UpsertFrom(in).LastErr; err != nil {
		return nil, err
	}
	self.module = m
	return self.module, err
}

func (self *Endpoint) FindTarget(target node.PathSlice, remote node.Node, local node.Node) node.Node {
	return &node.Extend{
		Node: local,
		OnExtend: func(e *node.Extend, sel *node.Selection, meta meta.MetaList, localChild node.Node) (node.Node, error) {
			if sel.Path().Equal(target.Tail) && meta.GetIdent() == target.Tail.Meta().GetIdent() {
				// Found target node, not get out of the picture
				return node.ConfigNode(remote, localChild), nil
			}
			return self.FindTarget(target, remote, localChild), nil
		},
		/* OnAction */
	}
}

func (self *Endpoint) Proxy(target node.PathSlice) (node.Node, error) {
	var path string
	if ! target.Empty() {
		path = target.String()
	}
	url := self.EndpointAddress + "restconf/" + path + "?content=nonconfig"
	c2.Info.Println(url)
	remote, err := self.DownloadJson(url)
	if err != nil {
		return nil, err
	}
	local := node.StoreNode(self.Store)
	if ! target.Empty() {
		return self.FindTarget(target, remote, local), nil
	}
	return node.ConfigNode(remote, local), nil
}

func (self *Endpoint) DownloadJson(url string) (node.Node, error) {
	var req *http.Request
	var err error
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := http.DefaultClient
	resp, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	defer resp.Body.Close()
	return node.NewJsonReader(resp.Body).Node(), nil
}
