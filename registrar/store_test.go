package registrar

import (
	"testing"
	"github.com/c2g/meta"
	"github.com/c2g/meta/yang"
	"github.com/c2g/node"
	"strings"
	"bytes"
)

var storeTestModule = `
module m {
  namespace "";
  prefix "";
  revision 0;
  container a {
    leaf b {
    	type string;
    }
    list c {
    	key "y";
    	leaf y {
    		type string;
    	}
    	container z {
    		type int32;
    	}
    }
  }
}
`

func TestStoreLoad(t *testing.T) {
	test := `a/b:string:X
a/c=y/y:string:Y
a/c=y/z:int32:99
`
	f := FileStorageHandler{}
	b := node.NewBufferStore()
	err := f.LoadFromReader(strings.NewReader(test), b)
	if err != nil {
		t.Error(err)
	}
	if len(b.Values) != 3 {
		t.Error(len(b.Values))
	}
	actualStr := b.Values["a/b"].Str
	if actualStr != "X" {
		t.Error(actualStr)
	}
	actualInt := b.Values["a/c=y/z"].Int
	if actualInt != 99 {
		t.Error(actualInt)
	}
}

func TestStoreSave(t *testing.T) {
	f := FileStorageHandler{}
	b := node.NewBufferStore()
	b.Values["a/b"] = &node.Value{
		Str:"X",
		Type: meta.NewDataType(nil, "string"),
	}
	b.Values["a/c=y/z"] = &node.Value{
		Int:77,
		Type: meta.NewDataType(nil, "int32"),
	}
	var actual bytes.Buffer
	err := f.SaveToWriter(&actual, b)
	if err != nil {
		t.Error(err)
	}
	expected := `a/b:string:X
a/c=y/z:int32:77
`
	if actual.String() != expected {
		t.Error("\n" + actual.String())
	}
}

func YangFromString(s string) *meta.Module {
	m, err := yang.LoadModuleCustomImport(s, nil)
	if err != nil {
		panic(err)
	}
	return m
}