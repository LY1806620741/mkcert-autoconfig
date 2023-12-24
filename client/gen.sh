#!bin/bash
current_path=$(cd "$(dirname $0)";pwd)
cd $current_path
mkdir dist -p
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/certClient_linux_x64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -o dist/certClient_windows_x64.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/certClient_mac_x64
cp index.html dist/index.html
cd ../
mkdir dist -p
mv -u $current_path/dist/* dist