#!/bin/bash
set -eu

protoc -I proto proto/*.proto --go_out=plugins=grpc:proto
