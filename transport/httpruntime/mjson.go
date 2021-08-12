package httpruntime

import (
	"io"

	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	// removed gogo support as it is incompatible with protobuf-v2
	return m.Unmarshaler.NewDecoder(r).Decode(dst)
}

func (m MarshalerPbJSON) Marshal(w io.Writer, src interface{}) error {
	// removed gogo support as it is incompatible with protobuf-v2
	return m.Marshaler.NewEncoder(w).Encode(src)
}
