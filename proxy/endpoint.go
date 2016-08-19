package proxy

import (
	"github.com/c2stack/c2g/browse"
	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/node"
	"io"
	"net/http"
	"fmt"
	"bytes"
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
	if err = node.SelectModule(self.YangPath, m, false).Root().UpsertFrom(in).LastErr; err != nil {
		return nil, err
	}
	self.Meta = m
	return self.Meta, err
}

// Navigate to given node
func navigate(target string, n node.Node) node.Node {
	e := &node.MyNode{}
	checkTarget := func(current string) (node.Node) {
		if target == current {
			// Found target node, now get out of the picture
			return n
		}

		return e
	}
	e.OnSelect = func(r node.ContainerRequest) (node.Node, error) {
		return checkTarget(fmt.Sprint(r.Selection.Path.StringNoModule(), "/", r.Meta.GetIdent())), nil
	}
	e.OnNext = func(r node.ListRequest) (node.Node, []*node.Value, error) {
		return checkTarget(fmt.Sprint(r.Selection.Path.StringNoModule(), "=", node.EncodeKey(r.Key))), r.Key, nil
	}
	return e
}

func (self *Endpoint) handleRequest(target node.PathSlice) (node.Node, error) {
	self.basePath = target.Head
	var path string
	if !target.Empty() {
		path = target.String()
	}
	operational, err := self.getRequest("restconf/" + path + "?content=nonconfig")
	if err != nil && ! c2.IsNotFoundErr(err) {
		return nil, err
	}
	tx, createTxErr := self.TxSource.ConfigStore(self)
	if createTxErr != nil {
		return nil, createTxErr
	}
	config, beginTxErr := tx.ConfigNode(path)
	if beginTxErr != nil && ! c2.IsNotFoundErr(beginTxErr) {
		return nil, beginTxErr
	}
	proxy := &proxy{
		stripPathPrefix: target.Head.String(),
		onCommit:        tx.SaveConfig,
		onRequest:       self.request,
	}
	prxy := proxy.proxy(config, operational)
	if len(path) == 0 {
		return prxy, nil
	}
	return navigate(target.Tail.StringNoModule(), prxy), nil
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

func (self *Endpoint) pushConfig() error {
	tx, createTxErr := self.TxSource.ConfigStore(self)
	if createTxErr != nil {
		return createTxErr
	}
	localConfig, beginTxErr := tx.ConfigNode("")
	if beginTxErr != nil && ! c2.IsNotFoundErr(beginTxErr) {
		return beginTxErr
	}
	var payload bytes.Buffer
	payloadNode := node.NewJsonWriter(&payload).Node()
	if err := node.NewBrowser2(self.Meta, localConfig).Root().InsertInto(payloadNode).LastErr; err != nil {
		return err
	}
	if _, err := self.request("PUT", "restconf/", &payload); err != nil {
		return err
	}
	return nil
}

func (self *Endpoint) pullConfig() error {
	remoteConfig, err := self.getRequest("restconf/?content=config")
	if err != nil {
		return err
	}
	tx, createTxErr := self.TxSource.ConfigStore(self)
	if createTxErr != nil {
		return createTxErr
	}
	localConfig, beginTxErr := tx.ConfigNode("")
	if beginTxErr != nil && ! c2.IsNotFoundErr(beginTxErr) {
		return beginTxErr
	}
	if err := node.NewBrowser2(self.Meta, localConfig).Root().UpsertFrom(remoteConfig).LastErr; err != nil {
		return err
	}
	return tx.SaveConfig()
}

func (self *Endpoint) getRequest(url string) (node.Node, error) {
	body, err := self.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return node.NewJsonReader(body).Node(), nil
}
