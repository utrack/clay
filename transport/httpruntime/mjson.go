package httpruntime

import (
	"io"
	"io/ioutil"
	"reflect"

	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	v := reflect.ValueOf(dst)

	// try to extract the message from under the pointers
	// and determine if protomsg is registered under gogo registry
	// see #35
	var isGogo bool
	for {
		if kind := v.Kind(); kind != reflect.Ptr && kind != reflect.Interface {
			break
		}
		vv := v.Interface()
		pm, ok := vv.(proto.Message)
		if ok {
			isGogo = gogoproto.MessageName(pm) != ""
			dst = vv
			break
		}
		v = v.Elem()
	}

	if pm, ok := dst.(proto.Message); ok {
		if isGogo {
			return m.GogoUnmarshaler.Unmarshal(r, pm)
		}
	}

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if isArray(body) {
		return m.GogoUnmarshaler.UnmarshalArray(body, dst)
	}

	return m.Unmarshaler.Unmarshal(body, dst)
}

func (m MarshalerPbJSON) Marshal(w io.Writer, src interface{}) error {
	if pm, ok := src.(proto.Message); ok {
		if gogoproto.MessageName(pm) != "" {
			return m.GogoMarshaler.Marshal(w, pm)
		}
	}
	if in := tryToMakeArrayWithData(src); in != nil {
		return m.GogoMarshaler.MarshalArray(w, in)
	}
	return m.Marshaler.NewEncoder(w).Encode(src)
}

type arrayWithData []interface{}
func tryToMakeArrayWithData(in interface{}) arrayWithData {
	switch reflect.TypeOf(in).Kind() {
	default:
		return nil
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(in)
		b := make(arrayWithData, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			b = append(b, s.Index(i).Interface())
		}
		return b
	}
}

func isArray(body []byte) bool {
	if len(body) > 0 && body[0] == '[' && body[len(body)-1] == ']' {
		return true
	}
	return false
}