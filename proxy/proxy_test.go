package proxy

import (
	"github.com/c2g/c2"
	"github.com/c2g/meta/yang"
	"github.com/c2g/node"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestProxy(t *testing.T) {
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
		t.Fatal(err)
	}
	operational := node.NewJsonReader(strings.NewReader(`{"a":{"aa":"a.aa"}}`)).Node()
	config := node.MapNode(map[string]interface{}{
		"b": map[string]interface{}{
			"bb": "b.bb",
		}})

	tests := []struct{
		edit string
		method string
		url string
		find string
	} {
		{
			edit : 	`{"c":{"cc":"c.cc"}}`,
			method : "POST",
			url : "restconf/test",
			find : "",
		},
		{
			edit : 	`{"bd":{"bda":{"bdaa":"bd.bda.bdaa"}}}`,
			method : "PUT",
			url : "restconf/test/b",
			find : "b",
		},
	}
	for i, test := range tests {
		t.Logf("Test #%d", i)
		edit := node.NewJsonReader(strings.NewReader(test.edit)).Node()
		var requestCalled, commitCalled bool
		proxy := &proxy{
			onRequest: func(method string, url string, payload io.Reader) (io.ReadCloser, error) {
				if requestCalled {
					t.Error("Request called multiple times")
				}
				requestCalled = true
				if err := c2.CheckEqual(test.method, method); err != nil {
					t.Error(err)
				}
				if err := c2.CheckEqual(test.url, url); err != nil {
					t.Error(err)
				}
				editPayload, _ := ioutil.ReadAll(payload)
				if err := c2.CheckEqual(test.edit, string(editPayload)); err != nil {
					t.Error(err)
				}
				return nil, nil
			},
			onCommit: func() error {
				commitCalled = true
				return nil
			},
		}
		n := proxy.proxy(config, operational)
		s := node.NewBrowser2(m, n).Root().Selector().Find(test.find)
		if err := s.InsertFrom(edit).LastErr; err != nil {
			t.Error(err)
		}
		if !commitCalled {
			t.Error("Commit never called")
		}
		if !requestCalled {
			t.Error("Request never called")
		}

	}
}
