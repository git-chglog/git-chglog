package main

const templateTagNameAnchor = "<a name=\"{{ .Tag.Name }}\"></a>\n"

// TemplateBuilder ...
type TemplateBuilder interface {
	Builder
}

// TemplateBuilderFactory ...
type TemplateBuilderFactory = func(string) TemplateBuilder

func templateBuilderFactory(template string) TemplateBuilder {
	switch template {
	case tplKeepAChangelog.display:
		return NewKACTemplateBuilder()
	default:
		return NewCustomTemplateBuilder()
	}
}
