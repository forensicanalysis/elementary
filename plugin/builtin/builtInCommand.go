package builtin

import (
	"bytes"
	"io/ioutil"

	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var Filter = &pluginlib.Parameter{Name: "filter", Description: "filter processed events", Type: pluginlib.StringArray, Required: false}

func fileToReader(store *forensicstore.ForensicStore, exportPath gjson.Result) (*bytes.Reader, error) {
	file, teardown, err := store.LoadFile(exportPath.String())
	if err != nil {
		return nil, err
	}
	defer teardown()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

/*
func setupWriterPlugin(p plugin.Plugin) (*forensicstore.ForensicStore, *output.MainWriter, func() error, error) {
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return nil, nil, nil, err
	}

	out := output.New(
		p.Parameter().StringValue("output"),
		p.Parameter().StringValue("format"),
	)
	if p.Parameter().BoolValue("add-to-store") {
		out.Add(plugin.NewForensicStoreOutput(store))
	}
	return store, out, func() error {
		out.WriteFooter()
		return teardown()
	}, nil
}
*/

func getForensicStore(p pluginlib.Plugin) (*forensicstore.ForensicStore, func() error, error) {
	path := p.Parameter().StringValue("forensicstore")
	return forensicstore.Open(path)
}
