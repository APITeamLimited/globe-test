# Makes a binary smaller by using upx
go build -ldflags="-s -w" -o worker_function main.go
# Use upx to shrink the binary
upx --best --lzma worker_function
