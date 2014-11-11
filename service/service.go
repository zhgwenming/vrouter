package service

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
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

// A TCP service will publish to outer network
type Service struct {
	Name     string
	Host     string // to be scheduled on
	Addr     string // VIP or local address
	Port     string
	Backends []*Backend
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
	return buf.Bytes()
}

func (s *Service) UnMarshal(buf []byte) error {
	buffer := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(buffer)
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
