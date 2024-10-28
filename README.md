# Simple Chat
`by Javo`

Simple chat running on TCP, implemented in Go.


## Install

`sh 
go mod download
`

## Run 
### As Server 

`sh 
go run main.go --mode server --ipAddress <your_ip_address> 
`
### As Client

`sh 
go run main.go --mode client --ipAddress <server_ip_address:server_port> 
go run main.go --ipAddress <server_ip_address:server_port> 
`
