package registrar

import (
	"bufio"
	"bytes"
	"github.com/c2g/node"
	"github.com/c2g/c2"
	"fmt"
	"io"
	"os"
	"github.com/c2g/meta"
	"sort"
)

// This implementation is not very performant on medium to large nodesets
type FileStorageHandler struct {
	Id  string
	Dir string
}

func (self *FileStorageHandler) Manage() node.Node {
	return node.MarshalContainer(self)
}

func (self *FileStorageHandler) Filename() string {
	return fmt.Sprintf("%s/%s.json", self.Dir, self.Id)
}

func (self *FileStorageHandler) LoadFromDisk(buff *node.BufferStore) error {
	f, err := os.OpenFile(self.Filename(), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return self.LoadFromReader(f, buff)
}

func (self *FileStorageHandler) LoadFromReader(ior io.Reader, buff *node.BufferStore) error {
	var val *node.Value
	var key string
	r := bufio.NewReader(ior)
	line, isprefix, err := r.ReadLine()
	for line != nil && err == nil {
		if isprefix {
			return c2.NewErr("Value too big for file reading buffer")
		}
		if key, val, err = self.Decode(line); err != nil {
			return err
		}
		buff.Values[key] = val
		line, isprefix, err = r.ReadLine()
	}
	return nil
}

func (self *FileStorageHandler) SaveToDisk(buff *node.BufferStore) error {
	f, err := os.OpenFile(self.Filename(), os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return self.SaveToWriter(f, buff)
}

func (self *FileStorageHandler) SaveToWriter(iow io.Writer, store *node.BufferStore) (err error) {
	w := bufio.NewWriter(iow)

	// TODO: Lock value store, - invalid on multi-thread access,
	keys := make([]string, len(store.Values))
	i := 0
	for k, _ := range store.Values {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := store.Values[k]
		if err = self.Encode(w, k, v); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func (self *FileStorageHandler) Decode(line []byte) (key string, v *node.Value, err error) {
	colon1 := bytes.IndexRune(line, ':')
	if colon1 < 0 {
		return "", nil, c2.NewErr("Missing colon " + string(line))
	}
	key = string(line[:colon1])
	formatAndValue := line[colon1+1:]
	colon2 := bytes.IndexRune(formatAndValue, ':')
	if colon2 < 0 {
		return "", nil, c2.NewErr("Missing second colon " + string(line))
	}
	format := string(formatAndValue[:colon2])
	v = &node.Value{
		Type: meta.NewDataType(nil, format),
	}
	vstr := string(formatAndValue[colon2+1:])
	if err = v.CoerseStrValue(vstr); err != nil {
		return "", nil, err
	}
	return
}

func (self *FileStorageHandler) Encode(w *bufio.Writer, k string, v *node.Value) (err error) {
	if _, err = w.WriteString(k); err != nil {
		return err
	}
	if _, err = w.WriteRune(':'); err != nil {
		return err
	}

	if _, err = w.WriteString(v.Type.Format().String()); err != nil {
		return err
	}

	if _, err = w.WriteRune(':'); err != nil {
		return err
	}

	// TODO: Encode string to remove CRLF
	if _, err = w.WriteString(v.String()); err != nil {
		return err
	}

	if _, err = w.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

func (self *FileStorageHandler) Store() node.Store {
	store := node.NewBufferStore()
	store.OnLoad = self.LoadFromDisk
	store.OnSave = self.SaveToDisk
	return store
}
