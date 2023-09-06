#!/bin/bash
#SERVERMAIN=./tokenServer/main.go
CLIENTMAIN=./tokenClient/main.go



#Running Client

go run $CLIENTMAIN -create -id 1 -host localhost -port 4200

go run $CLIENTMAIN -write -id 1 -name sai -low 0 -mid 10 -high 50 -host localhost -port 4200

go run $CLIENTMAIN -read -id 1 -host localhost -port 4200

go run $CLIENTMAIN -create -id 2 -host localhost -port 4200

go run $CLIENTMAIN -write -id 2 -name teja -low 5 -mid 10 -high 50 -host localhost -port 4200 &
go run $CLIENTMAIN -write -id 2 -name peddi -low 5 -mid 20 -high 100 -host localhost -port 4200 &
go run $CLIENTMAIN -read -id 1 -host localhost -port 4200


go run $CLIENTMAIN -drop -id 2 -host localhost -port 4200



