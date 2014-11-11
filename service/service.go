package service

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrNotImplemented = errors.New("Function not implemented")

// Service Backend
type Backend struct {
	Addr string
	Port string
}

func NewBackend(str string) (*Backend, error) {
	fields := strings.Split(str, ":")
	if len(fields) != 2 {
		err := fmt.Errorf("wrong format of backend: %s", str)
		return nil, err
	}

	host := fields[0]
	port := fields[1]
	return &Backend{host, port}, nil
}

func (back *Backend) String() string {
	return back.Addr + ":" + back.Port
}

// A TCP service will publish to outer network
type Service struct {
	Name       string
	Active     bool
	CreateTime time.Time
	Host       string // to be scheduled on
	Addr       string // VIP or local address
	Port       string
	Backends   []*Backend
}

func NewService() *Service {
	srv := new(Service)
	tgt := make([]*Backend, 0, 4)
	srv.Backends = tgt
	return srv
}

func (s *Service) AddBackend(b *Backend) {
	s.Backends = append(s.Backends, b)
}

func (s *Service) Marshal() []byte {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(s)

	src := buf.Bytes()

	ascii := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(ascii, src)
	return ascii
}

func (s *Service) UnMarshal(str string) error {
	src, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(src)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(s)
}

// BindIP binds the Addr to a physical interface
func (s *Service) bindIP() error {
	return ErrNotImplemented
}

func (s *Service) Start() error {
	if err := s.bindIP(); err != nil {
		return err
	}
	return ErrNotImplemented
}

func (s *Service) removeIP() error {
	return ErrNotImplemented
}

func (s *Service) Stop() error {
	if err := s.removeIP(); err != nil {
		return err
	}
	return ErrNotImplemented
}

type LBProxy struct {
}
