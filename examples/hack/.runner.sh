#!/bin/bash
dir=$(dirname $0)
APP_PASSWORD="my-shiny-password" GOWORK=off go run ${dir}/../main.go $@
