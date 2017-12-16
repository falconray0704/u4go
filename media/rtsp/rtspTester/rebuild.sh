#!/bin/bash

rm -rf rtspTester rtspTesterArm rtspTester32.exe rtspTester64.exe logs

go build rtspTester.go

#GOOS=windows GOARCH=amd64 go build -o rtspTester64.exe rtspTester.go
#GOOS=windows GOARCH=386 go build -o rtspTester32.exe rtspTester.go
#GOOS=linux GOARCH=arm64 go build -o rtspTesterArm rtspTester.go


