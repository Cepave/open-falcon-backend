#!/bin/bash

echo -n "Building DB Patch program for OSX ... "
GOOS=darwin GOARCH=386 go build -o dbpatch-osx-32 ../go/
echo -e "Finish\n"
echo -n "Building DB Patch program for Windows ... "
GOOS=windows GOARCH=386 go build -o dbpatch-windows-32 ../go/
echo -e "Finish\n"
echo -n "Building DB Patch program for Linux ... "
GOOS=linux GOARCH=386 go build -o dbpatch-linux-32 ../go/
echo -e "Finish\n"
