#!/bin/bash
set -ex
cd $(dirname $0)

go run ./*.go

