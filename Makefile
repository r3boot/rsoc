TARGET = rsoc

BUILD_DIR = ./build
PREFIX = /usr/local

CI_DIR = ./ci
COVERAGE_DIR = ${CI_DIR}/coverage

all: ${BUILD_DIR} dependencies ${BUILD_DIR}/${TARGET}
tests: ${COVERAGE_DIR} test coverage

${BUILD_DIR}:
	mkdir -p "${BUILD_DIR}"

dependencies:
	go get -v ./...

${BUILD_DIR}/${TARGET}:
	go build -v -o ${BUILD_DIR}/${TARGET} \
		cmd/${TARGET}/${TARGET}.go

${COVERAGE_DIR}:
	mkdir -p "${COVERAGE_DIR}"

test:
	go test ./...

coverage: ${COVERAGE_DIR}
	${CI_DIR}/run_tests.sh

codecov:
	find ${COVERAGE_DIR} -name *.out -exec cat {} \; > coverage.txt

install:
	install -o root -g root -m 0755 ${BUILD_DIR}/${TARGET} ${PREFIX}/bin/${TARGET}

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rf "${BUILD_DIR}" || true
	[[ -d "${COVERAGE_DIR}" ]] && rm -rf "${COVERAGE_DIR}" || true
