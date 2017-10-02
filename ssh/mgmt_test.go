package ssh

import (
	"testing"

	"github.com/c2stack/c2g/nodes"

	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/node"
)

func TestManage(t *testing.T) {
	s := NewService()
	model := yang.RequireModule(&meta.FileStreamSource{Root: "../yang"}, "ssh")
	b := node.NewBrowser(model, Manage(s))
	config := `{
		"address" : "127.0.0.1:2202",
		"hostKeyFiles" : ["testdata/server_key_rsa"],
		"authorizedKeysFile" : "testdata/authorized_keys"
	}`
	if err := b.Root().InsertFrom(nodes.ReadJSON(config)).LastErr; err != nil {
		t.Fatal(err)
	}
}
