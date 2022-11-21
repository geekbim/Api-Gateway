# Api Gateway - KRAKEND
Api gateway for https://github.com/geekbim/Go-Microservices

#### Build Plugin Krakend

1. Prerequisites -> OS Linux & go version 1.16.4
2. Build login.so -> cd plugins/login && go build -buildmode=plugin -o login.so ./plugin

##### How To Add New Endpoint
1. Edit file env/krakend.json
2. Run generate script -> ./env/generate.sh
3. Commit and push 2 files env/krakend.json & dist/krakend.json