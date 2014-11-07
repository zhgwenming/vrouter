package service

import (
	"bytes"
	"encoding/gob"
)

// Service Instance
type Instance struct {
	Addr string
	Port string
}

// A TCP service will publish to outer network
type Service struct {
	Name    string
	Addr    string
	Port    string
	Targets []*Instance
}

func NewService() *Service {
	srv := new(Service)
	tgt := make([]*Instance, 0, 4)
	srv.Targets = tgt
	return srv
}

func (s *Service) Marshal() []byte {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(s)
	return buf.Bytes()
}

func (s *Service) UnMarshal(buf []byte) {
	buffer := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(s)
}

type LBProxy struct {
}
