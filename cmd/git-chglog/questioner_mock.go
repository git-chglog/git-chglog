package main

type mockQuestionerImpl struct {
	ReturnAsk func() (*Answer, error)
}

func (m *mockQuestionerImpl) Ask() (*Answer, error) {
	return m.ReturnAsk()
}
