package genhandler

import "testing"

func Test_goTypeName(t *testing.T) {
	tests := []struct {
		in string
		want string
	}{
		{"test", "Test"},
		{"testVar", "TestVar"},
		{"test.testVar", "test.TestVar"},
		{"test.testVar.nestedVar", "test.TestVar.NestedVar"},
	}
	for _, tc := range tests {
		if got := goTypeName(tc.in); got != tc.want {
			t.Errorf("goTypeName(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
