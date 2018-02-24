package main

type mockFileSystem struct {
	ReturnExists    func(string) bool
	ReturnMkdirP    func(string) error
	ReturnCreate    func(string) (File, error)
	ReturnWriteFile func(string, []byte) error
}

func (m *mockFileSystem) Exists(path string) bool {
	return m.ReturnExists(path)
}

func (m *mockFileSystem) MkdirP(path string) error {
	return m.ReturnMkdirP(path)
}

func (m *mockFileSystem) Create(name string) (File, error) {
	return m.ReturnCreate(name)
}

func (m *mockFileSystem) WriteFile(path string, content []byte) error {
	return m.ReturnWriteFile(path, content)
}

type mockFile struct {
	File
	ReturnWrite func([]byte) (int, error)
}

func (m *mockFile) Write(b []byte) (int, error) {
	return m.ReturnWrite(b)
}
