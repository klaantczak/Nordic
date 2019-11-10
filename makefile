# Makefile should be placed into the root folder of the golang project.
export GOPATH=$(CURDIR)

VERSION_FILE = src/hps/version.gen.go

# Contains project source files, including generated
SRC = $(shell find src -name "*.go") $(VERSION_FILE)

# Just the current date in format 2016-05-08 01:29
BUILD_DATE = $(shell date +%Y-%m-%d\ %H:%M)

# Contains full git revision with optional `-modified` suffix indicating that
# the files in the working folder are modified. For example:
# c99c8bcbd457386a04e1eb50efc468ad474abb87-modified
BUILD_ID = $(shell if [ -d .git ]; then git describe --abbrev=40 --dirty=-modified --always; fi)

.PHONY: all deps version compile test
.PHONY: nordic32-compile nordic32-test nordic32
.PHONY: tasks
.PHONY: nordic32lf-compile nordic32lf
.PHONY: softdev-compile softdev
.PHONY: package pull-hps push-hps

all: test nordic32-test tasks softdev-compile

deps:
	go get github.com/gonum/matrix/mat64
	go get github.com/stretchr/testify/assert
	go get gonum.org/v1/plot/...

# On every compilation $(VERSION_FILE) is updated with current version and date.
# Constants, defined in this file can be saved to log file or displayed when
# the version key is specified in command arguments. This file is ignored by
# source control engine. If BUILD_ID is empty (GIT is not available), the
# file will not be updated.
version:
ifneq ($(BUILD_ID),)
	@echo "BUILD: $(BUILD_ID) $(BUILD_DATE)"
	@echo "package hps" > $(VERSION_FILE)
	@echo "const (" >> $(VERSION_FILE)
	@echo "  BUILD_ID = \"$(BUILD_ID)\"" >> $(VERSION_FILE)
	@echo "  BUILD_DATE = \"$(BUILD_DATE)\"" >> $(VERSION_FILE)
	@echo ")" >> $(VERSION_FILE)
endif
ifeq ($(BUILD_ID),)
	@echo "BUILD: $(BUILD_DATE)"
	@echo "package hps" > $(VERSION_FILE)
	@echo "const (" >> $(VERSION_FILE)
	@echo "  BUILD_ID = \"DATE\"" >> $(VERSION_FILE)
	@echo "  BUILD_DATE = \"$(BUILD_DATE)\"" >> $(VERSION_FILE)
	@echo ")" >> $(VERSION_FILE)
endif

bin/hpscmd: version $(SRC)
	go fmt hps/... hpstest/... hpscmd/...
	go install hpscmd

compile: bin/hpscmd

test: compile
	go test hpstest -timeout 2s

bin/nordic32: $(SRC)
	go fmt nordic32/... nordic32test/...
	go install nordic32

nordic32-compile: bin/nordic32

nordic32-test: nordic32-compile
	go test nordic32test/... -timeout 2s

nordic32: nordic32-test

bin/tasks: $(SRC)
	go fmt tasks/...
	go install tasks

tasks: bin/tasks

nordic32lf-compile:
	go fmt nordic32lf/...
	go install -buildmode=c-archive nordic32lf
	gcc -o bin/nordic32lfbin -I pkg/darwin_amd64 src/nordic32lfbin/main.c pkg/darwin_amd64/nordic32lf.a

nordic32lf: nordic32lf-compile

bin/softdev: $(SRC)
	go fmt softdev/...
	go install softdev

softdev-compile: bin/softdev

softdev: softdev-compile

# Creates redistributable HPS package, omitting .* folders and files, compiled
# binaries, logs and legacy files. Inclues all the third-party libraries.
# The created package is ignored by source control.
package: version
	COPYFILE_DISABLE=1 tar czf hps.tgz \
		--exclude ".*" \
		--exclude bin \
		--exclude hps.tgz \
		--exclude logs \
		--exclude pkg \
		* 

# Updates HPS engine from the public package. The public package should
# be available at ../hps.
pull-hps:
	rm -rf src/hps src/hpscmd src/hpstest
	cp -r ../hps/src/hps src
	cp -r ../hps/src/hpscmd src
	cp -r ../hps/src/hpstest src

push-hps:
	rm -rf ../hps/src/hps ../hps/src/hpscmd ../hps/src/hpstest
	cp -r src/hps ../hps/src
	cp -r src/hpscmd ../hps/src
	cp -r src/hpstest ../hps/src
