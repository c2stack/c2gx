package proxy

import (
	"github.com/c2g/c2"
	"github.com/c2g/meta"
	"github.com/c2g/node"
	"io"
	"net/http"
	"github.com/c2g/browse"
)

type Endpoint struct {
	YangPath        meta.StreamSource
	Id              string
	Module          string
	Meta            *meta.Module
	EndpointAddress string
	TxSource        ConfigStoreSource
	ClientSource    browse.ClientSource
	basePath        *node.Path
}

func (self *Endpoint) Schema() (*meta.Module, error) {
	if self.Meta != nil {
		return self.Meta, nil
	}
	in, err := self.getRequest("meta/?noexpand")
	if err != nil {
		return nil, err
	}
	m := &meta.Module{}
	if err = node.SelectModule(self.YangPath, m, false).Root().Selector().UpsertFrom(in).LastErr; err != nil {
		return nil, err
	}
	self.Meta = m
	return self.Meta, err
}

func (self *Endpoint) FindTarget(target node.PathSlice, n node.Node) node.Node {
	e := &node.Extend{
		OnExtend: func(e *node.Extend, sel *node.Selection, meta meta.MetaList, localChild node.Node) (node.Node, error) {
			if sel.Path().Equal(target.Tail) && meta.GetIdent() == target.Tail.Meta().GetIdent() {
				// Found target node, now get out of the picture
				return n, nil
			}

			// recursive
			return e, nil
		},
	}
	e.Node = e
	return e
}

func (self *Endpoint) handleRequest(target node.PathSlice) (node.Node, error) {
	self.basePath = target.Head
	var path string
	if !target.Empty() {
		path = target.String()
	}
	operational, err := self.getRequest("restconf/" + path + "?content=nonconfig")
	if err != nil {
		return nil, err
	}
	tx, createTxErr := self.TxSource.ConfigStore(self)
	if createTxErr != nil {
		return nil, createTxErr
	}
	config, beginTxErr := tx.ConfigNode(path)
	if beginTxErr != nil {
		return nil, beginTxErr
	}
	proxy := &proxy{
		stripPathPrefix: target.Head.String(),
		onCommit:        tx.SaveConfig,
		onRequest:       self.request,
	}
	return proxy.proxy(config, operational), nil
}

func (self *Endpoint) request(method string, url string, payload io.Reader) (io.ReadCloser, error) {
	var req *http.Request
	var err error
	fullUrl := self.EndpointAddress + url
	if req, err = http.NewRequest(method, fullUrl, payload); err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := self.ClientSource.GetHttpClient()
	c2.Info.Printf("%s %s", method, fullUrl)
	resp, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	return resp.Body, nil
}

func (self *Endpoint) getRequest(url string) (node.Node, error) {
	body, err := self.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return node.NewJsonReader(body).Node(), nil
}
