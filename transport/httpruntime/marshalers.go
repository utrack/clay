package httpruntime

import (
	"io"
	"net/http"
	"strings"

	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type Marshaler interface {
	ContentType() string
	Unmarshal(io.Reader, proto.Message) error
	Marshal(io.Writer, proto.Message) error
}

var marshalDict = map[string]Marshaler{
	"application/json": MarshalerPbJSON{
		Marshaler:       &jsonpb.Marshaler{},
		Unmarshaler:     &jsonpb.Unmarshaler{},
		GogoMarshaler:   &gogojsonpb.Marshaler{},
		GogoUnmarshaler: &gogojsonpb.Unmarshaler{},
	},
}

// OverrideMarshaler replaces marshaler for given content-type.
func OverrideMarshaler(contentType string, m Marshaler) {
	marshalDict[strings.ToLower(contentType)] = m
}

// MarshalerForRequest returns marshalers for inbound and outbound bodies.
func MarshalerForRequest(r *http.Request) (Marshaler, Marshaler) {
	inbound := marshalerOrDefault(r.Header.Get("Content-Type"))
	outbound := marshalerOrDefault(r.Header.Get("Accept"))
	return inbound, outbound
}

func marshalerOrDefault(t string) Marshaler {
	sepIdx := strings.IndexAny(t, ";,")
	// TODO we're not negotiating really. Account the q= param and additional
	// options
	if sepIdx > 0 {
		t = t[:sepIdx]
	}
	t = strings.ToLower(t)

	if m, ok := marshalDict[t]; ok {
		return m
	}
	return marshalDict[MarshalerPbJSON{}.ContentType()]
}

type MarshalerPbJSON struct {
	Marshaler       *jsonpb.Marshaler
	Unmarshaler     *jsonpb.Unmarshaler
	GogoMarshaler   *gogojsonpb.Marshaler
	GogoUnmarshaler *gogojsonpb.Unmarshaler
}

func (MarshalerPbJSON) ContentType() string {
	return "application/json"
}

func (m MarshalerPbJSON) Unmarshal(r io.Reader, dst proto.Message) error {
	if gogoproto.MessageName(dst) != "" {
		return m.GogoUnmarshaler.Unmarshal(r, dst)
	}
	return m.Unmarshaler.Unmarshal(r, dst)
}

func (m MarshalerPbJSON) Marshal(w io.Writer, src proto.Message) error {
	if gogoproto.MessageName(src) != "" {
		return m.GogoMarshaler.Marshal(w, src)
	}
	return m.Marshaler.Marshal(w, src)
}
