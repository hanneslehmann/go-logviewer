#!/bin/bash
./preprocess -input=logviewer.go -output=tmp.go -config=local.config
go build -o logviewer_local tmp.go
