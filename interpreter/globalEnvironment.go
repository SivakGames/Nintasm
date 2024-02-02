package interpreter

type Environment struct {
	name   string
	record map[string]string
	parent string
}

func NewEnvironment(envName string) Environment {
	return Environment{
		name: envName,
	}
}

func (e *Environment) Lookup(name string) string {
	return e.resolve(name)
}

func (e *Environment) resolve(name string) string {

	value, ok := e.record[name]
	if ok {
		return value
	}
	return ""
}
