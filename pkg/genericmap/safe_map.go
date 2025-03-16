package genericmap

import "sync"

type SafeMap[K comparable, V any] struct {
	mutex sync.RWMutex
	data  *Map[K, V]
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: NewMap[K, V](),
	}
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data.Set(key, value)
}

func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Get(key)
}

func (m *SafeMap[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data.Delete(key)
}

func (m *SafeMap[K, V]) Has(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Has(key)
}

func (m *SafeMap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Size()
}

func (m *SafeMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data.Clear()
}

func (m *SafeMap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Keys()
}

func (m *SafeMap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Values()
}

func (m *SafeMap[K, V]) ForEach(fn func(key K, value V) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	m.data.ForEach(fn)
}

func (m *SafeMap[K, V]) GetOrSet(key K, defaultValue V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data.GetOrSet(key, defaultValue)
}

func (m *SafeMap[K, V]) GetAndDelete(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data.GetAndDelete(key)
}

func (m *SafeMap[K, V]) Items() []Pair[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data.Items()
}
