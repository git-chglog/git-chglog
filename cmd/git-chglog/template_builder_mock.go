package main

type mockTemplateBuilderImpl struct {
	ReturnBuild func(*Answer) (string, error)
}

func (m *mockTemplateBuilderImpl) Build(ans *Answer) (string, error) {
	return m.ReturnBuild(ans)
}
