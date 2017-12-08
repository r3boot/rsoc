TARGET = rsoc

BUILD_DIR = ./build
BUILD_BRANCH = build

all: prepare dependencies ${BUILD_DIR}/${TARGET}

prepare:
	git branch -l | grep -q ${BUILD_BRANCH} \
		&& git branch -D ${BUILD_BRANCH} \
		|| true
	git checkout -b "${BUILD_BRANCH}"
	mkdir -p ${BUILD_DIR}

dependencies:
	go get -v ./...

${BUILD_DIR}/${TARGET}:
	go build -v -o ${BUILD_DIR}/${TARGET} \
		cmd/${TARGET}/${TARGET}.go

clean:
	git checkout master
	git branch -D "${BUILD_BRANCH}"
