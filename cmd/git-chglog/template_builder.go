package main

// TemplateBuilder ...
type TemplateBuilder interface {
	Builder
}

// TemplateBuilderFactory ...
type TemplateBuilderFactory = func(string) TemplateBuilder

func templateBuilderFactory(template string) TemplateBuilder {
	switch template {
	case tplKeepAChangelog:
		return NewKACTemplateBuilder()
	default:
		return NewCustomTemplateBuilder()
	}
}
