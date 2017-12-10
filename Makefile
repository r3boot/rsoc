TARGET = rsoc

BUILD_DIR = ./build

CI_DIR = ./ci
COVERAGE_DIR = ${CI_DIR}/coverage

all: ${BUILD_DIR} ${BUILD_DIR}/${TARGET}

${BUILD_DIR}:
	mkdir -p "${BUILD_DIR}"

${BUILD_DIR}/${TARGET}:
	go build -v -o ${BUILD_DIR}/${TARGET} \
		cmd/${TARGET}/${TARGET}.go

${COVERAGE_DIR}:
	mkdir -p "${COVERAGE_DIR}"

test:
	go test ./...

coverage: ${COVERAGE_DIR}
	${CI_DIR}/run_tests.sh

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rf "${BUILD_DIR}" || true
	[[ -d "${COVERAGE_DIR}" ]] && rm -rf "${COVERAGE_DIR}" || true
