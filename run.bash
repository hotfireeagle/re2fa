# !/bin/bash
go run main.go | (cd ./frontend && npm run dev)