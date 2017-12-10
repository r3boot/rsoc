TARGET = rsoc

BUILD_DIR = ./build
TEST_DIR = ./test_results
CI_DIR = ./ci

all: ${BUILD_DIR} ${BUILD_DIR}/${TARGET}

${BUILD_DIR}:
	mkdir -p "${BUILD_DIR}"

${BUILD_DIR}/${TARGET}:
	go build -v -o ${BUILD_DIR}/${TARGET} \
		cmd/${TARGET}/${TARGET}.go

${TEST_DIR}:
	mkdir -p "${TEST_DIR}"

unittests: ${TEST_DIR}
	${CI_DIR}/run_tests.sh

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rf "${BUILD_DIR}" || true
