#!/bin/bash

DIR="$( cd "$( dirname "$0" )" && pwd )"
go run $DIR/cmd/main.go ${@}
