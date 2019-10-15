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

func (data Data) GetMulti(keys ...string) []interface{} {
	var result = make([]interface{}, 0, len(keys))
	for _, key := range keys {
		v, _ := data[key]
		result = append(result, v)
	}
	return result
}

func (data Data) Pop(key string) interface{} {
	if v, ok := data[key]; ok {
		delete(data, key)
		return v
	}
	return nil
}
