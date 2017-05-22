package httpruntime

import (
	"io"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type Marshaler interface {
	ContentType() string
	Unmarshal(io.Reader, proto.Message) error
	Marshal(io.Writer, proto.Message) error
}

var defaultMarshaler = MarshalerPbJSON{Marshaler: &jsonpb.Marshaler{}}

var marshalDict = map[string]Marshaler{
	"application/json": defaultMarshaler,
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
	sepIdx := strings.Index(t, ";")
	// TODO we're not negotiating really. Account the q= param and additional
	// options
	if sepIdx > 0 {
		t = t[:sepIdx]
	}
	t = strings.ToLower(t)

	if m, ok := marshalDict[t]; ok {
		return m
	}
	return defaultMarshaler
}

type MarshalerPbJSON struct {
	Marshaler *jsonpb.Marshaler
}

func (MarshalerPbJSON) ContentType() string {
	return "application/json"
}

func (MarshalerPbJSON) Unmarshal(r io.Reader, dst proto.Message) error {
	return jsonpb.Unmarshal(r, dst)
}

func (m MarshalerPbJSON) Marshal(w io.Writer, src proto.Message) error {
	return m.Marshaler.Marshal(w, src)
}
