#!/bin/bash

CWD=$(pwd)
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Generate sentimentanalyzer enum.
cd "${SCRIPT_DIR}/../pkg/sentimentanalyzer"
go generate

cd "$CWD"
