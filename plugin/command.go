package plugin

type Annotation string

const (
	Importer Annotation = "importer"
	Exporter Annotation = "exporter"
	Di       Annotation = "di"
)

func HasAnnotation(cmd Plugin, annotation Annotation) bool {
	for _, a := range cmd.Annotations() {
		if a == annotation {
			return true
		}
	}
	return false
}

type Plugin interface {
	Name() string
	Short() string
	Parameter() ParameterList
	Run(Plugin) error
	Annotations() []Annotation
}
