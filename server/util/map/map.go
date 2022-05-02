package Map

type StandardMap map[any]any

func New() StandardMap {
	return make(map[any]any)
}

func (this StandardMap) ContainsKey(key any) bool {
	return this[key] != nil
}

func (this StandardMap) ContainsValue(value any) bool {
	return this.GetKey(value) != nil
}

func (this StandardMap) Put(key any, value any) StandardMap {
	this[key] = value
	return this
}

func (this StandardMap) Get(key any) any {
	return this[key]
}

func (this StandardMap) GetKey(value any) any {
	for key, currentValue := range this {
		if value == currentValue {
			return key
		}
	}
	return nil
}

func (this StandardMap) GetOrDefault(key any, defaultValue any) any {
	if this.ContainsKey(key) {
		return this[key]
	}
	return defaultValue
}

func (this StandardMap) Remove(key string) any {
	data := this[key]
	this[key] = nil
	return data
}

func (this StandardMap) Size() int {
	return len(this)
}

func (this StandardMap) IsEmpty() bool {
	return len(this) != 0
}

func (this StandardMap) Each(iterator func(key, value any)) {
	for key, value := range this {
		iterator(key, value)
	}
}

func (this StandardMap) Clear() StandardMap {
	for key, _ := range this {
		this[key] = nil
	}
	return this
}
