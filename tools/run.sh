#
#   Auto run project 
#
go mod download

cd tools
go run init_localenv.go
cd ..

air --build.cmd "go build -o ./bin/api" --build.bin "./bin/api"