package protobuf

//go:generate protoc --go_out=. sample.proto

import (
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func TestProto2(t *testing.T) {
	p0 := &Person{
		Name: proto.String("Alice"),
		Age: proto.Int(30),
	}
	b, err := proto.Marshal(p0)
	if err != nil {
		t.Fatal(err)
	}
	p1 := &Person{}
	if err := proto.Unmarshal(b, p1); err != nil {
		t.Fatal(err)
	}
	text := proto.MarshalTextString(p1)
	want := "name: \"Alice\"\nage: 30\n"
	if text != want {
		t.Errorf("MarshalTextString() = %q; Want %q", text, want)
	}
}

func TestProto2JSON(t *testing.T) {
	p := &Person{
		Name: proto.String("Alice"),
		Age: proto.Int(30),
	}
	m := jsonpb.Marshaler{}
	json, err := m.MarshalToString(p)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\"name\":\"Alice\",\"age\":30}"
	if json != want {
		t.Errorf("got %q; want %q", json, want)
	}
}
