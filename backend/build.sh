#!/usr/bin/env bash

export APP_NAME="go-parquet-read"
export ENV="dev"
export PARQUET_FILE_DIRECTORT="/Users/viadmin/Downloads/full-stack-engineer/Assignment-Full-Stack-Apica/dev"

go install
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
 
echo "Doing some cleaning ..."
go clean
echo "Done."

# Check if required commands are available
command -v go >/dev/null 2>&1 || { echo >&2 "Go is required but it's not installed. Aborting."; exit 1; }
command -v gofmt >/dev/null 2>&1 || { echo >&2 "Gofmt is required but it's not installed. Aborting."; exit 1; }

echo "Running go format ..."
gofmt -w .
echo "Done."

echo "Running go build ..."
go build -race
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
echo "Done."

if [ $? == 0 ]; then
	echo "## Starting service ##"
    ./go-parquet-read
fi