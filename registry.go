package sii

import (
	"reflect"
	"sync"
)

var (
	blockClass       = map[string]reflect.Type{}
	blockClassLock   = new(sync.RWMutex)
	defaultBlockType = reflect.TypeOf(RawBlock{})
)

func RegisterBlock(b Block) {
	blockClassLock.Lock()
	defer blockClassLock.Unlock()

	blockClass[b.Class()] = reflect.TypeOf(b).Elem()
}

func getBlockInstance(t string) Block {
	blockClassLock.RLock()
	defer blockClassLock.RUnlock()

	if rt, ok := blockClass[t]; ok {
		v := reflect.New(rt).Interface()
		if b, ok := v.(Block); ok {
			return b
		}
	}

	return reflect.New(defaultBlockType).Interface().(Block)
}
