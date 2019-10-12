package event

type Data map[string]interface{}

func (data Data) Set(key string, v interface{}) {
	data[key] = v
}

func (data Data) Get(key string) interface{} {
	if v, ok := data[key]; ok {
		return v
	}
	return nil
}

func (data Data) Pop(key string) interface{} {
	if v, ok := data[key]; ok {
		delete(data, key)
		return v
	}
	return nil
}
