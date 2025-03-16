package genericmap

type Map[K comparable, V any] struct {
	data map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]V),
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	m.data[key] = value
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	val, exists := m.data[key]
	return val, exists
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.data, key)
}

func (m *Map[K, V]) Has(key K) bool {
	_, exists := m.data[key]
	return exists
}

func (m *Map[K, V]) Size() int {
	return len(m.data)
}

func (m *Map[K, V]) Clear() {
	m.data = make(map[K]V)
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, v)
	}
	return values
}

func (m *Map[K, V]) ForEach(fn func(key K, value V) bool) {
	for k, v := range m.data {
		if !fn(k, v) {
			break
		}
	}
}

func (m *Map[K, V]) GetOrSet(key K, defaultValue V) V {
	if val, exists := m.data[key]; exists {
		return val
	}
	m.data[key] = defaultValue
	return defaultValue
}

func (m *Map[K, V]) GetAndDelete(key K) (V, bool) {
	val, exists := m.data[key]
	if exists {
		delete(m.data, key)
	}
	return val, exists
}

func (m *Map[K, V]) Items() []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m.data))
	for k, v := range m.data {
		pairs = append(pairs, Pair[K, V]{k, v})
	}
	return pairs
}
