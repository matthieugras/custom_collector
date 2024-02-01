#!/bin/bash
set -e
go generate .
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" .
mv custom-collector.exe collect_files.exe
zip collect_files.zip collect_files.exe
bytes=$(cat collect_files.zip | base64 -w0)
sed -f - script_template.txt > collect_files.ps1 << EOF
s~<EXE_PLACEHOLDER>~${bytes}~
EOF
