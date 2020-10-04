#!/bin/bash
# Simple shell script to build application for linux
ENV env GOOS=linux GOARCH=amd64 go build -o bin cmd/server.go
