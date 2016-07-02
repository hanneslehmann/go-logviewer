#!/bin/bash
./preprocess -input=logviewer.go -output=tmp.go -config=local.config
go run tmp.go
