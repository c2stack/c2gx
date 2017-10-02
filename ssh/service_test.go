package ssh

import "testing"

func TestService(t *testing.T) {
	s := NewService()
	o := s.Options()
	o.Address = "0.0.0.0:2202"
	o.HostKeyFiles = []string{
		"testdata/server_key_rsa",
	}
	o.AuthorizedKeysFile = "testdata/authorized_keys"
	if err := s.Apply(o); err != nil {
		t.Fatal(err)
	}
}
