package host

import (
	"testing"
	"time"

	"github.com/c2stack/c2g/c2"

	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/nodes"

	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/node"
)

func TestHostNode(t *testing.T) {
	c2.DebugLog(true)
	model := yang.RequireModule(&meta.FileStreamSource{Root: "../yang"}, "host")
	b := node.NewBrowser(model, Manage())
	metrics := &Metrics{SampleRate: 100}
	go metrics.Start()
	<-time.After(100 * time.Millisecond)
	defer metrics.closer()

	s, err := nodes.WriteJSON(b.Root())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
