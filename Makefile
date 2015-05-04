NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=/bin/echo -e
endif

all: 
	@$(ECHO) "$(OK_COLOR)==> Building$(NO_COLOR)"
	go get -v ./...
	go test -v ./...
	go build

bin: 
	@$(ECHO) "$(OK_COLOR)==> Building$(NO_COLOR)"
	go build

test: 
	@$(ECHO) "$(OK_COLOR)==> Testing$(NO_COLOR)"
	go get -v ./...
	go test -v ./...

clean:
	@rm -rf dist/ packer-post-processor-redis

format:
	go fmt ./...

dist:
	@$(ECHO) "$(OK_COLOR)==> Building Packages...$(NO_COLOR)"
	@gox -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 freebsd/386 freebsd/amd64 openbsd/386 openbsd/amd64 windows/386 windows/amd64 netbsd/386 netbsd/amd64"
	@mv packer-post-processor-redis_darwin_386 packer-post-processor-redis; tar cvfz packer-post-processor-redis.darwin-i386.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_darwin_amd64 packer-post-processor-redis; tar cvfz packer-post-processor-redis.darwin-amd64.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_freebsd_386 packer-post-processor-redis; tar cvfz packer-post-processor-redis.freebsd-i386.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_freebsd_amd64 packer-post-processor-redis; tar cvfz packer-post-processor-redis.freebsd-amd64.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_linux_386 packer-post-processor-redis; tar cvfz packer-post-processor-redis.linux-i386.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_linux_amd64 packer-post-processor-redis; tar cvfz packer-post-processor-redis.linux-amd64.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_netbsd_386 packer-post-processor-redis; tar cvfz packer-post-processor-redis.netbsd-i386.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_netbsd_amd64 packer-post-processor-redis; tar cvfz packer-post-processor-redis.netbsd-amd64.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_openbsd_386 packer-post-processor-redis; tar cvfz packer-post-processor-redis.openbsd-i386.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_openbsd_amd64 packer-post-processor-redis; tar cvfz packer-post-processor-redis.openbsd-amd64.tar.gz packer-post-processor-redis; rm packer-post-processor-redis
	@mv packer-post-processor-redis_windows_386.exe packer-post-processor-redis.exe; zip packer-post-processor-redis.windows-i386.zip packer-post-processor-redis.exe; rm packer-post-processor-redis.exe
	@mv packer-post-processor-redis_windows_amd64.exe packer-post-processor-redis.exe; zip packer-post-processor-redis.windows-amd64.zip packer-post-processor-redis.exe; rm packer-post-processor-redis.exe
	@mkdir -p dist/
	@mv packer-post-processor-redis* dist/.

.PHONY: all clean deps format test updatedeps