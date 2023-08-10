package pEtcd

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var Default *Etcd

func StoreDefault(etcd *Etcd) {
	Default = etcd
}

func String(name string) (string, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return "", false
	}
	return string(bs), true
}

func PromiseString(name, defVal string) string {
	s, ok := Default.String(name)
	if !ok {
		return defVal
	}
	return s
}

func Int(name string) (int, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return int(i64), true
}

func PromiseInt(name string, defVal int) int {
	v, ok := Default.Int(name)
	if !ok {
		return defVal
	}
	return v
}

func Int8(name string) (int8, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i8, err := strconv.ParseInt(string(bs), 10, 8)
	if err != nil {
		return 0, false
	}
	return int8(i8), true
}

func PromiseInt8(name string, defVal int8) int8 {
	v, ok := Default.Int8(name)
	if !ok {
		return defVal
	}
	return v
}

func Int16(name string) (int16, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i16, err := strconv.ParseInt(string(bs), 10, 16)
	if err != nil {
		return 0, false
	}
	return int16(i16), true
}

func PromiseInt16(name string, defVal int16) int16 {
	v, ok := Default.Int16(name)
	if !ok {
		return defVal
	}
	return v
}

func Int32(name string) (int32, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i32, err := strconv.ParseInt(string(bs), 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(i32), true
}

func PromiseInt32(name string, defVal int32) int32 {
	v, ok := Default.Int32(name)
	if !ok {
		return defVal
	}
	return v
}

func Int64(name string) (int64, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return i64, true
}

func PromiseInt64(name string, defVal int64) int64 {
	v, ok := Default.Int64(name)
	if !ok {
		return defVal
	}
	return v
}

func Uint(name string) (uint, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(i64), true
}

func PromiseUint(name string, defVal uint) uint {
	v, ok := Default.Uint(name)
	if !ok {
		return defVal
	}
	return v
}

func Uint8(name string) (uint8, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i8, err := strconv.ParseUint(string(bs), 10, 8)
	if err != nil {
		return 0, false
	}
	return uint8(i8), true
}

func PromiseUint8(name string, defVal uint8) uint8 {
	v, ok := Default.Uint8(name)
	if !ok {
		return defVal
	}
	return v
}

func Uint16(name string) (uint16, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i16, err := strconv.ParseUint(string(bs), 10, 16)
	if err != nil {
		return 0, false
	}
	return uint16(i16), true
}

func PromiseUint16(name string, defVal uint16) uint16 {
	v, ok := Default.Uint16(name)
	if !ok {
		return defVal
	}
	return v
}

func Uint32(name string) (uint32, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i32, err := strconv.ParseUint(string(bs), 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(i32), true
}

func PromiseUint32(name string, defVal uint32) uint32 {
	v, ok := Default.Uint32(name)
	if !ok {
		return defVal
	}
	return v
}

func Uint64(name string) (uint64, bool) {
	bs, ok := Default.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return i64, true
}

func PromiseUint64(name string, defVal uint64) uint64 {
	v, ok := Default.Uint64(name)
	if !ok {
		return defVal
	}
	return v
}

// Bool Numeric value only 0 will be treated as false, other numeric
func Bool(name string) (bool, ok bool) {
	bs, ok := Default.load(name)
	if !ok {
		return false, false
	}
	s := string(bs)

	si, err := strconv.Atoi(s)
	if err == nil {
		if si == 0 {
			return false, true
		}
		return true, true
	}

	s = strings.ToLower(s)
	if s == "true" || s == "t" {
		return true, true
	}

	return false, true
}

func PromiseBool(name string, defVal bool) bool {
	v, ok := Default.Bool(name)
	if !ok {
		return defVal
	}
	return v
}

func JSON(name string, v any) error {
	bs, ok := Default.load(name)
	if !ok {
		return errors.Errorf("no such key found")
	}
	return json.Unmarshal(bs, v)
}
