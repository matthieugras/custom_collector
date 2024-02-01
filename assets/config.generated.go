package assets

import (
	"github.com/forensicanalysis/artifactcollector/collection"
	"github.com/spf13/afero"
)

var Config = &collection.Configuration{Artifacts: []string{"WindowsAllEventLogs", "BrowserHistory", "BrowserCache", "BrowserHistory"}, User: false, Case: "", OutputDir: "", FS: afero.Fs(nil)}
