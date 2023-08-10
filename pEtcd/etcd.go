package pEtcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Kind uint

const (
	EvtInit Kind = iota
	EvtPut
	EvtDelete
)

type WatchCallback func(kind Kind, key string, value []byte)

type Etcd struct {
	client  *clientv3.Client
	kvstore sync.Map
}

func New(cfg *Config) (*Etcd, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return &Etcd{
		client:  client,
		kvstore: sync.Map{},
	}, nil
}

func (r *Etcd) Close() error {
	return r.client.Close()
}

func (r *Etcd) Client() *clientv3.Client {
	return r.client
}

func (r *Etcd) Set(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return r.client.Put(context.TODO(), key, val, opts...)
}

func (r *Etcd) set(k string, v []byte) {
	r.kvstore.Store(k, v)
}

func (r *Etcd) remove(k string) {
	r.kvstore.Delete(k)
}

func (r *Etcd) Watch(prefix string, watchCallback ...WatchCallback) error {
	var fn WatchCallback
	if len(watchCallback) > 0 {
		fn = watchCallback[0]
	}
	reply, err := r.client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range reply.Kvs {
		r.set(string(kv.Key), kv.Value)
		if fn != nil {
			fn(EvtInit, string(kv.Key), kv.Value)
		}
	}
	go r.watch(prefix, fn)
	return nil
}

func (r *Etcd) watch(prefix string, callback WatchCallback) {
	reply := r.client.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for rep := range reply {
		if rep.Err() != nil {
			continue
		}
		for _, evt := range rep.Events {
			key := string(evt.Kv.Key)
			switch evt.Type {
			case mvccpb.PUT:
				r.set(key, evt.Kv.Value)
				if r.hasJSONCache(key) {
					r.removeJSONCache(key)
				}
				if callback != nil {
					callback(EvtPut, key, evt.Kv.Value)
				}
			case mvccpb.DELETE:
				r.remove(key)
				if r.hasJSONCache(key) {
					r.removeJSONCache(key)
				}
				if callback != nil {
					callback(EvtDelete, key, []byte{})
				}
			}
		}
	}
}

func (r *Etcd) load(name string) ([]byte, bool) {
	b, ok := r.kvstore.Load(name)
	if !ok {
		return nil, false
	}
	bs := b.([]byte)
	return bs, true
}

func (r *Etcd) String(name string) (string, bool) {
	bs, ok := r.load(name)
	if !ok {
		return "", false
	}
	return string(bs), true
}

func (r *Etcd) PromiseString(name, defVal string) string {
	s, ok := r.String(name)
	if !ok {
		return defVal
	}
	return s
}

func (r *Etcd) Int(name string) (int, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return int(i64), true
}

func (r *Etcd) PromiseInt(name string, defVal int) int {
	v, ok := r.Int(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Int8(name string) (int8, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i8, err := strconv.ParseInt(string(bs), 10, 8)
	if err != nil {
		return 0, false
	}
	return int8(i8), true
}

func (r *Etcd) PromiseInt8(name string, defVal int8) int8 {
	v, ok := r.Int8(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Int16(name string) (int16, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i16, err := strconv.ParseInt(string(bs), 10, 16)
	if err != nil {
		return 0, false
	}
	return int16(i16), true
}

func (r *Etcd) PromiseInt16(name string, defVal int16) int16 {
	v, ok := r.Int16(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Int32(name string) (int32, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i32, err := strconv.ParseInt(string(bs), 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(i32), true
}

func (r *Etcd) PromiseInt32(name string, defVal int32) int32 {
	v, ok := r.Int32(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Int64(name string) (int64, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return i64, true
}

func (r *Etcd) PromiseInt64(name string, defVal int64) int64 {
	v, ok := r.Int64(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Uint(name string) (uint, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(i64), true
}

func (r *Etcd) PromiseUint(name string, defVal uint) uint {
	v, ok := r.Uint(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Uint8(name string) (uint8, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i8, err := strconv.ParseUint(string(bs), 10, 8)
	if err != nil {
		return 0, false
	}
	return uint8(i8), true
}

func (r *Etcd) PromiseUint8(name string, defVal uint8) uint8 {
	v, ok := r.Uint8(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Uint16(name string) (uint16, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i16, err := strconv.ParseUint(string(bs), 10, 16)
	if err != nil {
		return 0, false
	}
	return uint16(i16), true
}

func (r *Etcd) PromiseUint16(name string, defVal uint16) uint16 {
	v, ok := r.Uint16(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Uint32(name string) (uint32, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i32, err := strconv.ParseUint(string(bs), 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(i32), true
}

func (r *Etcd) PromiseUint32(name string, defVal uint32) uint32 {
	v, ok := r.Uint32(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) Uint64(name string) (uint64, bool) {
	bs, ok := r.load(name)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return 0, false
	}
	return i64, true
}

func (r *Etcd) PromiseUint64(name string, defVal uint64) uint64 {
	v, ok := r.Uint64(name)
	if !ok {
		return defVal
	}
	return v
}

// Bool Numeric value only 0 will be treated as false, other numeric
func (r *Etcd) Bool(name string) (bool, ok bool) {
	bs, ok := r.load(name)
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

func (r *Etcd) PromiseBool(name string, defVal bool) bool {
	v, ok := r.Bool(name)
	if !ok {
		return defVal
	}
	return v
}

func (r *Etcd) JSON(name string, v any) error {
	key := jsonKey(name)
	if jsonVal, ok := r.kvstore.Load(key); ok {
		return r.json(jsonVal, v)
	}
	bs, ok := r.load(name)
	if !ok {
		return errors.Errorf("no such key found")
	}
	err := json.Unmarshal(bs, v)
	if err != nil {
		return err
	}
	r.kvstore.Store(key, v)
	return nil
}

func (r *Etcd) json(s any, t any) error {
	tTyp := reflect.TypeOf(t)
	sTyp := reflect.TypeOf(s)
	tVal := reflect.ValueOf(t).Elem()
	sVal := reflect.ValueOf(s).Elem()

	if tTyp.Kind() != reflect.Ptr {
		return errors.New("invalid kind of target value")
	}
	if tTyp != sTyp {
		return errors.New("invalid typ")
	}
	for i := 0; i < tVal.NumField(); i++ {
		// 在源结构体中查询有数据结构体中相同属性和类型的字段，有则修改其值
		// name := sTypeOfT.Field(i).Name
		f := tVal.Type().Field(i)
		name := f.Name
		targetFiled := sVal.FieldByName(name)
		if targetFiled.IsValid() && f.Type == targetFiled.Type() {
			tVal.FieldByName(name).Set(reflect.ValueOf(sVal.FieldByName(name).Interface()))
			tVal.Field(i).Set(reflect.ValueOf(sVal.FieldByName(name).Interface()))
		}
	}
	return nil
}

func (r *Etcd) hasJSONCache(key string) bool {
	_, ok := r.kvstore.Load(jsonKey(key))
	return ok
}

func (r *Etcd) removeJSONCache(key string) {
	r.kvstore.Delete(jsonKey(key))
}

func jsonKey(key string) string {
	return fmt.Sprintf("__KVSTORE_%s_JSON__", key)
}
