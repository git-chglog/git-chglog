package main

type mockFileSystem struct {
	ReturnMkdirP func(string) error
	ReturnCreate func(string) (File, error)
}

func (m *mockFileSystem) MkdirP(path string) error {
	return m.ReturnMkdirP(path)
}

func (m *mockFileSystem) Create(name string) (File, error) {
	return m.ReturnCreate(name)
}

type mockFile struct {
	File
	ReturnWrite func([]byte) (int, error)
}

func (m *mockFile) Write(b []byte) (int, error) {
	return m.ReturnWrite(b)
}
