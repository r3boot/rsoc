TARGET = rsoc

BUILD_DIR = ./build
BUILD_BRANCH = build

all: prepare dependencies ${BUILD_DIR}/${TARGET}

prepare:
	mkdir -p "${BUILD_DIR}"

dependencies:
	go get -v ./...

${BUILD_DIR}/${TARGET}:
	go build -v -o ${BUILD_DIR}/${TARGET} \
		cmd/${TARGET}/${TARGET}.go

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rf "${BUILD_DIR}" || true
