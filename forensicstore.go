package elementary

import (
	"log"

	"github.com/forensicanalysis/forensicstore"
)

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
