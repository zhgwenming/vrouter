package service

import (
	"bytes"
	"encoding/gob"
)

// Service Backend
type Backend struct {
	Addr string
	Port string
}

// A TCP service will publish to outer network
type Service struct {
	Name     string
	Host     string
	Addr     string
	Port     string
	Backends []*Backend
}

func NewService() *Service {
	srv := new(Service)
	tgt := make([]*Backend, 0, 4)
	srv.Backends = tgt
	return srv
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

type LBProxy struct {
}
