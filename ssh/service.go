package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/c2stack/c2g/c2"

	gssh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Service struct {
	config         *gssh.ServerConfig
	authorizedKeys map[string]struct{}
	stop           func()
	options        Options
}

type Options struct {
	Address            string
	HostKeyFiles       []string
	AuthorizedKeysFile string
}

func NewService() *Service {
	s := &Service{}
	return s
}

func (self *Service) Options() Options {
	return self.options
}

func (self *Service) Apply(options Options) error {
	if self.stop != nil {
		self.stop()
	}
	self.config = &gssh.ServerConfig{
		PasswordCallback:  self.passwordAuth,
		PublicKeyCallback: self.keyAuth,
	}
	for _, k := range options.HostKeyFiles {
		privateBytes, err := ioutil.ReadFile(k)
		if err != nil {
			return c2.NewErr(fmt.Sprintf("Failed to load private key %s : %v", k, err))
		}
		private, err := gssh.ParsePrivateKey(privateBytes)
		if err != nil {
			return c2.NewErr(fmt.Sprintf("Failed to parse private key %s: %v", k, err))
		}
		self.config.AddHostKey(private)
	}

	authBytes, err := ioutil.ReadFile(options.AuthorizedKeysFile)
	if err != nil {
		return c2.NewErr(fmt.Sprintf("Failed to read authorized file %s : %v", options.AuthorizedKeysFile, err))
	}
	self.authorizedKeys = make(map[string]struct{})
	for len(authBytes) > 0 {
		pubKey, _, _, rest, err := gssh.ParseAuthorizedKey(authBytes)
		if err != nil {
			log.Fatal(err)
		}
		self.authorizedKeys[string(pubKey.Marshal())] = struct{}{}
		authBytes = rest
	}

	self.options = options
	go self.start()
	return nil
}

func (self *Service) passwordAuth(c gssh.ConnMetadata, pass []byte) (*gssh.Permissions, error) {
	return nil, nil
}

func (self *Service) keyAuth(c gssh.ConnMetadata, pubKey gssh.PublicKey) (*gssh.Permissions, error) {
	return nil, nil
}

func (self *Service) start() error {
	l, err := net.Listen("tcp", self.options.Address)
	if err != nil {
		return err
	}

	c, err := l.Accept()
	if err != nil {
		return err
	}

	_, chans, reqs, err := gssh.NewServerConn(c, self.config)
	if err != nil {
		return err
	}

	go gssh.DiscardRequests(reqs)

	for newChannel := range chans {
		fmt.Println("new channel")
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(gssh.UnknownChannelType, "unknown channel type")
		}
		ch, req, err := newChannel.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}
		go func(in <-chan *gssh.Request) {
			for req := range in {
				req.Reply(req.Type == "shell", nil)
			}
		}(req)
		term := terminal.NewTerminal(ch, "> ")
		go func() {
			defer ch.Close()
			for {
				line, err := term.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(line)
			}
		}()
	}
	return nil
}
