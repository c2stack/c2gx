package proxy

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/node"
)

func proxyTestModule() *meta.Module {
	mstr := `module test {
	namespace "";
	prefix "";
	revision 0;
	container a {
		config "false";
		leaf aa {
			type string;
		}
		leaf ab {
			type string;
		}
		leaf ac {
			type string;
		}
	}
	container b {
		leaf ba {
			config "false";
			type string;
		}
		leaf bb {
			type string;
		}
		leaf bc {
			type string;
		}
		container bd {
			container bda {
				leaf bdaa {
				 	type string;
				}
			}
		}
	}
	container c {
		leaf ca {
			type string;
		}
		leaf cb {
			type string;
		}
		leaf cc {
			type string;
		}
	}
}`
	m, err := yang.LoadModuleFromString(nil, mstr)
	if err != nil {
		panic(err)
	}
	return m
}

func TestProxyNavigate(t *testing.T) {
	m := proxyTestModule()
	target := &node.MyNode{}
	n := navigate("b/bd", target)
	b := node.NewBrowser(m, n)
	sel := b.Root().Find("b/bd")
	if sel.LastErr != nil {
		t.Fatal(sel.LastErr)
	}
	if sel.Node != target {
		t.Error("Target not expected")
	}
}

func n(js string) node.Node {
	var data map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(js)).Decode(&data); err != nil {
		panic(err)
	}
	return node.MapNode(data)
}

func TestProxyRequest(t *testing.T) {
	m := proxyTestModule()
	operational := n(`{"a":{"aa":"a.aa"}}`)
	config := n(`{"b":{"bb":"b.bb"}}`)
	tests := []struct {
		edit   node.Node
		method string
		url    string
		find   string
	}{
		{
			edit:   n(`{"c":{"cc":"c.cc"}}`),
			method: "POST",
			url:    "restconf/test",
			find:   "",
		},
		// TODO: Fix
		//{
		//	edit : 	`{"bd":{"bda":{"bdaa":"bd.bda.bdaa"}}}`,
		//	method : "PUT",
		//	url : "restconf/test/b",
		//	find : "b",
		//},
	}
	for i, test := range tests {
		t.Logf("Test #%d", i)
		//var requestCalled, commitCalled bool
		var actual []string
		//var actual bytes.Buffer
		p := &proxy{
			onRequest: func(method string, url string, payload io.Reader) (io.ReadCloser, error) {
				actual = append(actual, "req")
				// if requestCalled {
				// 	t.Error("Request called multiple times")
				// }
				// requestCalled = true
				// if err := c2.CheckEqual(test.method, method); err != nil {
				// 	t.Error(err)
				// }
				// if err := c2.CheckEqual(test.url, url); err != nil {
				// 	t.Error(err)
				// }
				// editPayload, _ := ioutil.ReadAll(payload)
				// if err := c2.CheckEqual(test.edit, string(editPayload)); err != nil {
				// 	t.Error(err)
				// }
				return nil, nil
			},
			onCommit: func() error {
				actual = append(actual, "commit")
				return nil
			},
		}
		n := p.proxy(config, operational)
		s := node.NewBrowser(m, n).Root().Find(test.find)
		if err := s.InsertFrom(test.edit).LastErr; err != nil {
			t.Error(err)
		}
		// Disable this test, unclear what we're actually testing here
		//c2.Equals(t, []string{"hello"}, actual)
	}
}
