export GOPATH=/root/go_fcgi
export GOBIN=/root/go_fcgi
go build -gcflags "-N -l" -race -o xx xx 
gdb xx