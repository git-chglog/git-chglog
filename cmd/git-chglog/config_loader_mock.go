package main

type mockConfigLoaderImpl struct {
	ReturnLoad func(string) (*Config, error)
}

func (m *mockConfigLoaderImpl) Load(path string) (*Config, error) {
	return m.ReturnLoad(path)
}
