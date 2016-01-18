.PHONY: .FORCE
GO=go
DOT=dot
GOYACC=$(GO) tool yacc

SRC_DIR = ./src
PROTO_INSTALL_FILE_DIR = ./src/code.google.com/p/goprotobuf/

all:
	$(GO) install DaemonServer

clean:
	rm -rf bin pkg
 
fmt:
	$(GO) fmt $(SRC_DIR)/...
	
#交叉编译：
#首先进入go源码目录
#cd /usr/local/go/src/
#执行sudo GOOS=linux GOARCH=amd64 ./make.bash  生成linux下的编译文件pkg
#以上步骤在go1.5中应该已不再需要，但还未测试
publish:
	GOOS=linux GOARCH=amd64 $(GO) build -o DaemonServer DaemonServer
	