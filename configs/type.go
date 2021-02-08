package configs

type Type struct {
	Map map[string]string
}

func (t *Type) List() map[string]string {
	return t.Map
}

func (t *Type) Value(key string) string {
	if value, ok := t.Map[key]; ok {
		return value
	}

	return ""
}
