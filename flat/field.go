package flat

import (
	"encoding"
	"reflect"
	"strconv"
	"time"
)

var _ Field = (*field)(nil)

type field struct {
	name string
	meta map[string]string

	tag   reflect.StructTag
	field reflect.Value
}

// Used by standard library flag package.
func (f *field) IsBoolFlag() bool {
	return f.field.Kind() == reflect.Bool
}

func (f *field) Name() string {
	return f.name
}

func (f *field) Meta() map[string]string {
	return f.meta
}

func (f *field) Tag(key string) (string, bool) {
	return f.tag.Lookup(key)
}

func (f *field) String() string {
	return f.tag.Get("default")
}

func (f *field) Get() interface{} {
	// return f.ptr
	return nil
}

var textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func (f *field) Set(value string) error {

	t := f.field.Type()

	if t.Implements(textUnmarshalerType) {
		return f.setUnmarshale([]byte(value))
	}

	switch f.field.Kind() {
	case reflect.String:
		return f.setString(value)
	case reflect.Bool:
		return f.setBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.String() == "time.Duration" {
			return f.setDuration(value)
		}
		return f.setInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.setUint(value)
	case reflect.Float32, reflect.Float64:
		return f.setFloat(value)

		// Soon case reflect.Map:
		// Soon case reflect.Slice:

		// Maybe case reflect.Func:
		// Maybe case reflect.Array:

		// Why? case reflect.Complex64:
		// Why? case reflect.Complex128:

		// Never case reflect.Chan:
		// Never case reflect.Interface:
		// Never case reflect.Ptr:
		// Never case reflect.Struct:
		// Never case reflect.UnsafePointer:
	}
	return nil
}

func (f *field) setUnmarshale(value []byte) error {

	if f.field.IsNil() {
		f.field.Set(reflect.New(f.field.Type().Elem()))
	}
	ut := f.field.MethodByName("UnmarshalText")

	err := ut.Call([]reflect.Value{reflect.ValueOf(value)})[0]

	if err.IsNil() {
		return nil
	}
	return err.Interface().(error)
}

func (f *field) setDuration(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	f.field.SetInt(int64(duration))
	return nil
}

func (f *field) setString(value string) error {
	f.field.SetString(value)
	return nil
}

func (f *field) setBool(value string) error {
	v, err := strconv.ParseBool(value)
	f.field.SetBool(v)
	return err
}

func (f *field) setInt(value string) error {
	v, err := strconv.ParseInt(value, 0, 64)
	f.field.SetInt(v)
	return err
}

func (f *field) setUint(value string) error {
	v, err := strconv.ParseUint(value, 0, 64)
	f.field.SetUint(v)
	return err
}

func (f *field) setFloat(value string) error {
	v, err := strconv.ParseFloat(value, 64)
	f.field.SetFloat(v)
	return err
}
