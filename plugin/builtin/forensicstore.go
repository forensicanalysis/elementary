package builtin

import (
	"log"

	"github.com/forensicanalysis/elementary/plugin/output"
	"github.com/forensicanalysis/forensicstore"
)

var _ output.Writer = &ForensicStoreOutput{}

type ForensicStoreOutput struct {
	store *forensicstore.ForensicStore
}

func NewForensicStoreOutput(store *forensicstore.ForensicStore) *ForensicStoreOutput {
	return &ForensicStoreOutput{store: store}
}

func (o *ForensicStoreOutput) WriteHeader([]string) {}

func (o *ForensicStoreOutput) WriteLine(element []byte) {
	_, err := o.store.Insert(element)
	if err != nil {
		log.Println(err, string(element))
	}
}

func (o *ForensicStoreOutput) WriteFooter() {}
