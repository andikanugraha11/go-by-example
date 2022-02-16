#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. greet\greetpb\greet.proto