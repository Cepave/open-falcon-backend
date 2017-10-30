#!/bin/bash

echo -n "Building DB Patch program for OSX ... "
GOOS=darwin GOARCH=386 go build -o dbpatch-osx-32 ../go/
GOOS=darwin GOARCH=amd64 go build -o dbpatch-osx-64 ../go/
echo -e "Finish\n"
echo -n "Building DB Patch program for Windows ... "
GOOS=windows GOARCH=386 go build -o dbpatch-windows-32 ../go/
GOOS=windows GOARCH=amd64 go build -o dbpatch-windows-64 ../go/
echo -e "Finish\n"
echo -n "Building DB Patch program for Linux ... "
GOOS=linux GOARCH=386 go build -o dbpatch-linux-32 ../go/
GOOS=linux GOARCH=amd64 go build -o dbpatch-linux-64 ../go/
echo -e "Finish\n"
