package tools

type JsonMap map[string]interface{}

func (m JsonMap) GetString(key string) (string, bool) {
	if value, ok := m[key]; ok {
		if v, ok := value.(string); ok {
			return v, true
		}
	}
	return "", false
}

func (m JsonMap) GetInt(key string) (int, bool) {
	if value, ok := m[key]; ok {
		if v, ok := value.(int); ok {
			return v, true
		}
	}
	return 0, false
}

func (m JsonMap) GetFloat(key string) (float64, bool) {
	if value, ok := m[key]; ok {
		if v, ok := value.(float64); ok {
			return v, true
		}
	}
	return 0.0, false
}

func (m JsonMap) GetListOfStrings(key string) ([]string, bool) {
	if value, ok := m[key]; ok {
		if list, ok := value.([]interface{}); ok {
			result := make([]string, 0)
			for _, value := range list {
				if str, ok := value.(string); ok {
					result = append(result, str)
					continue
				}
				return []string{}, false
			}
			return result, true
		}
	}
	return []string{}, false
}

func (m JsonMap) GetMap(key string) (JsonMap, bool) {
	if value, ok := m[key]; ok {
		if jmap, ok := value.(map[string]interface{}); ok {
			return JsonMap(jmap), true
		}
	}
	return JsonMap(make(map[string]interface{})), false
}
