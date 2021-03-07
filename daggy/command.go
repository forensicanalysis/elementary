package daggy

type Annotation string

const (
	Importer Annotation = "importer"
	Exporter            = "exporter"
	Di                  = "di"
)

func HasAnnotation(cmd Command, annotation Annotation) bool {
	for _, a := range cmd.Annotations() {
		if a == annotation {
			return true
		}
	}
	return false
}

type Command interface {
	Name() string
	Short() string
	Parameter() ParameterList
	Run(Command) error
	Annotations() []Annotation
}

