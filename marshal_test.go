package protohuman

import (
	"os"
	"testing"

	"github.com/amenzhinsky/vplay/encode/testdata"
	"github.com/golang/protobuf/proto"
)

func TestEncode(t *testing.T) {
	if err := Marshal(os.Stderr, &testdata.Test{
		Uint32F:     666,
		StringSlice: []string{"a", "b"},
		Inner:       &testdata.Test_Inner{Enabled: true},
		Kv:          map[string]int32{"four": 4, "five": 5},
		Buf:         []byte{1, 2, 3},
		Oneofer:     &testdata.Test_One{One: "one"},
	}); err != nil {
		t.Fatal(err)
	}

	proto.MarshalText()
}
