package protohuman

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

func Marshal(w io.Writer, msg proto.Message) error {
	return (&Marshaller{}).Marshal(w, msg)
}

func MarshalCompact(w io.Writer, msg proto.Message) error {
	return (&Marshaller{Compact: true}).Marshal(w, msg)
}

type Marshaller struct {
	Compact bool
	indent  int
}

func (m *Marshaller) Marshal(w io.Writer, msg proto.Message) error {
	wr := &writer{w: w}
	m.message(wr, msg, false)
	return wr.err
}

type writer struct {
	w   io.Writer
	err error
}

func (w *writer) write(s string) {
	if w.err != nil {
		return
	}
	if _, err := w.w.Write([]byte(s)); err != nil {
		w.err = err
	}
}

func (m *Marshaller) nl(w *writer) {
	if m.Compact {
		return
	}
	w.write("\n")
	w.write(strings.Repeat("\t", m.indent))
}

func (m *Marshaller) message(w *writer, v interface{}, isInterface bool) {
	if v == nil {
		w.write("nil")
		return
	}
	val := reflect.ValueOf(v)
	if val.IsNil() {
		w.write("nil")
		return
	}
	val = val.Elem()
	typ := val.Type()

	if isInterface {
		w.write("(")
	} else {
		w.write("{")
		m.indent++
		m.nl(w)
	}
	for i := 0; i < val.NumField(); i++ {
		name := fieldName(typ.Field(i))
		if name == "" {
			continue
		}
		if i != 0 {
			w.write(",")
			if !m.Compact {
				m.nl(w)
			}
		}
		w.write(name)
		w.write(":")
		if !m.Compact {
			w.write(" ")
		}
		m.value(w, val.Field(i))
	}
	if isInterface {
		w.write(")")
	} else {
		m.indent--
		m.nl(w)
		w.write("}")
	}
}

func fieldName(typ reflect.StructField) string {
	oneof, ok := typ.Tag.Lookup("protobuf_oneof")
	if ok {
		return oneof
	}
	tag, ok := typ.Tag.Lookup("protobuf")
	if !ok {
		return ""
	}
	for _, s := range strings.Split(tag, ",") {
		if !strings.HasPrefix(s, "name=") {
			continue
		}
		return s[5:]
	}
	return ""
}

func (m *Marshaller) value(w *writer, val reflect.Value) {
	switch val.Kind() {
	//case reflect.Map:
	case reflect.Slice:
		w.write("[")
		for i := 0; i < val.Len(); i++ {
			if i != 0 {
				w.write(",")
				if !m.Compact {
					w.write(" ")
				}
			}
			m.value(w, val.Index(i))
		}
		w.write("]")
	case reflect.String:
		w.write(fmt.Sprintf("%q", val.Interface()))
	case reflect.Ptr:
		m.message(w, val.Interface().(proto.Message), false)
	case reflect.Interface:
		m.message(w, val.Interface(), true)
	default:
		w.write(fmt.Sprintf("%v", val.Interface()))
	}
}
