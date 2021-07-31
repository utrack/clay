package genhandler

import (
	pbdescriptor "google.golang.org/protobuf/types/descriptorpb"
)

// copied from genswagger
// used to provide types for path generator
func primitiveTypeToGo(t pbdescriptor.FieldDescriptorProto_Type) (format string) {
	ret := ""
	switch t {
	case pbdescriptor.FieldDescriptorProto_TYPE_DOUBLE:
		ret = "float64"
	case pbdescriptor.FieldDescriptorProto_TYPE_FLOAT:
		ret = "float"
	case pbdescriptor.FieldDescriptorProto_TYPE_INT64:
		ret = "int64"
	case pbdescriptor.FieldDescriptorProto_TYPE_UINT64:
		ret = "uint64"
	case pbdescriptor.FieldDescriptorProto_TYPE_INT32:
		ret = "int32"
	case pbdescriptor.FieldDescriptorProto_TYPE_FIXED64:
		ret = "uint64"
	case pbdescriptor.FieldDescriptorProto_TYPE_FIXED32:
		ret = "int64"
	case pbdescriptor.FieldDescriptorProto_TYPE_BOOL:
		ret = "bool"
	case pbdescriptor.FieldDescriptorProto_TYPE_STRING:
		ret = "string"
	case pbdescriptor.FieldDescriptorProto_TYPE_BYTES:
		ret = "byte"
	case pbdescriptor.FieldDescriptorProto_TYPE_UINT32:
		ret = "uint32"
	case pbdescriptor.FieldDescriptorProto_TYPE_SFIXED32:
		ret = "int32"
	case pbdescriptor.FieldDescriptorProto_TYPE_SFIXED64:
		ret = "int64"
	case pbdescriptor.FieldDescriptorProto_TYPE_SINT32:
		ret = "int32"
	case pbdescriptor.FieldDescriptorProto_TYPE_SINT64:
		ret = "int64"
	default:
		ret = "string"
	}
	return ret
}
