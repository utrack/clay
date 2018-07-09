package httpruntime

import (
	"io"

	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"reflect"
)

var mpbjson = MarshalerPbJSON{
	Marshaler:       &runtime.JSONPb{},
	Unmarshaler:     &runtime.JSONPb{},
	GogoMarshaler:   &gogojsonpb.Marshaler{},
	GogoUnmarshaler: &gogojsonpb.Unmarshaler{},
}

// MarshalerPbJSON (un)marshals between JSON and proto.Messages.
// It supports both golang/pb and gogo/pb.
type MarshalerPbJSON struct {
	Marshaler       *runtime.JSONPb
	Unmarshaler     *runtime.JSONPb
	GogoMarshaler   *gogojsonpb.Marshaler
	GogoUnmarshaler *gogojsonpb.Unmarshaler
}

func (MarshalerPbJSON) ContentType() string {
	return "application/json"
}

func (m MarshalerPbJSON) Unmarshal(r io.Reader, dst interface{}) error {
	logrus.Error(reflect.TypeOf(dst).String())
	v := reflect.ValueOf(dst).Elem()
	if v.Kind() == reflect.Ptr { // pointer to pointer
		dst = v.Interface()
	}
	if pm, ok := dst.(proto.Message); ok {
		if gogoproto.MessageName(pm) != "" {
			logrus.Error("gogo")
			return m.GogoUnmarshaler.Unmarshal(r, pm)
		}
	}
	logrus.Error("go")
	return m.Unmarshaler.NewDecoder(r).Decode(dst)
}

func (m MarshalerPbJSON) Marshal(w io.Writer, src interface{}) error {
	if pm, ok := src.(proto.Message); ok {
		if gogoproto.MessageName(pm) != "" {
			return m.GogoMarshaler.Marshal(w, pm)
		}
	}
	return m.Marshaler.NewEncoder(w).Encode(src)
}
