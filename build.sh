go build -o ./bin/vnt-explorer main.go
go build -o ./bin/dbsync tools/dbsync/sync.go
go build -o ./bin/racer tools/racer/racer.go
go build -o ./bin/feixiaohao tools/feixiaohao/feixiaohao.go
# go build -o ./bin/mytoken tools/mytoken/mytoken.go
go build -o ./bin/nodemonitor tools/nodemonitor/node_monitor.go
go build -o ./bin/geniuser tools/racer/geniuser.go
