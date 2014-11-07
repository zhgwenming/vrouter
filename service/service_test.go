package service

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	orig := NewService()
	orig.Name = "test"
	orig.Addr = "127.0.0.1"
	orig.Port = "80"

	inst := new(Instance)
	inst.Addr = "172.16.1.1"
	inst.Port = "9199"

	orig.Targets = append(orig.Targets, inst)
	buf := orig.Marshal()

	dec := NewService()
	dec.UnMarshal(buf)

	if !reflect.DeepEqual(orig, dec) {
		t.Fatalf("error with unmarshaled value, orig: %#v, new: %#v", orig, dec)
	}
}
