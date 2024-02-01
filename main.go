// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package main

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/forensicanalysis/artifactcollector/collection"
	"github.com/forensicanalysis/artifactcollector/run"
	"github.com/forensicanalysis/artifactcollector/zipwrite"
	"github.com/forensicanalysis/artifactlib/goartifacts"
	"github.com/forensicanalysis/custom-collector/assets"
	split_writer "github.com/forensicanalysis/custom-collector/splitwriter"
)

//go:generate go install golang.org/x/tools/cmd/goimports@v0.1.7
//go:generate go install github.com/cugu/go-resources/cmd/resources@v0.3.0
//go:generate go install github.com/akavel/rsrc@v0.10.2
//go:generate go run scripts/yaml2go/main.go pack/ac.yaml pack/artifacts/*
//go:generate resources -package assets -output assets/bin.generated.go pack/bin/*
//go:generate rsrc -arch amd64 -manifest resources/artifactcollector.exe.manifest -ico resources/artifactcollector.ico -o resources/artifactcollector.syso
//go:generate rsrc -arch 386 -manifest resources/artifactcollector32.exe.manifest -ico resources/artifactcollector.ico -o resources/artifactcollector32.syso
//go:generate rsrc -arch amd64 -manifest resources/artifactcollector.exe.user.manifest -ico resources/artifactcollector.ico -o resources/artifactcollector.user.syso
//go:generate rsrc -arch 386 -manifest resources/artifactcollector32.exe.user.manifest -ico resources/artifactcollector.ico -o resources/artifactcollector32.user.syso

func main() {
	var maxsize uint = math.MaxUint
	if len(os.Args) > 1 {
		maxsize64, err := strconv.ParseUint(os.Args[1], 0, 64)
		if err != nil {
			log.Fatalf("Invalid parameter \"%s\", expected uint64", os.Args[0])
		}
		maxsize = uint(maxsize64)
	}
	f, err := split_writer.New("collection", maxsize)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var artifacts []goartifacts.ArtifactDefinition
	artifacts = append(artifacts, assets.Artifacts...)

	var fs = zipwrite.New(&f)
	defer fs.Close()
	if fs == nil {
		log.Fatal("Zipwrite is nil")
	}

	config := collection.Configuration{
		Artifacts: assets.Config.Artifacts,
		FS:        fs,
	}

	run.Run(&config, artifacts, nil)
}
