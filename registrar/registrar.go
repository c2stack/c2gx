package registrar

import (
	"fmt"
	"github.com/c2g/meta"
	"time"
)

type Registrar struct {
	Endpoints            map[string]*Endpoint
	SchemaInsertionPoint *meta.Choice
	StoreDir             string
}

func (self *Registrar) RegisterEndpoint(endpoint *Endpoint) error {
	if len(endpoint.Id) == 0 {

		// TODO: use real UUID
		uuid := fmt.Sprintf("%s.%d", endpoint.Module, time.Now().UnixNano())

		endpoint.Id = uuid
	}

	fs := &FileStorageHandler{Dir: self.StoreDir, Id: endpoint.Id}
	endpoint.Store = fs.Store()

	if module, err := endpoint.Schema(); err != nil {
		return err
	} else {
		kase := &meta.ChoiceCase{Ident: module.GetIdent()}
		self.SchemaInsertionPoint.AddMeta(kase)
		kase.AddMeta(module)
	}
	if self.Endpoints == nil {
		self.Endpoints = make(map[string]*Endpoint)
	}
	self.Endpoints[endpoint.Id] = endpoint
	return nil
}
