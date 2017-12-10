#!/usr/bin/env bash

BASE_DIR="$(dirname $(dirname $(readlink -f ${0})))"
OUTPUT_DIR="${BASE_DIR}/test_results"

if [[ ! -d "${OUTPUT_DIR}" ]]; then
    mkdir -p "${OUTPUT_DIR}"
fi

find ${BASE_DIR} -name *_test.go | xargs dirname | sort | uniq | while read SUBDIR; do
    TEST_NAME="$(basename ${SUBDIR})"
    COVERAGE_RESULTS="${OUTPUT_DIR}/${TEST_NAME}.out"
    COVERAGE_HTML="${OUTPUT_DIR}/${TEST_NAME}.html"

    echo -e "\n>>> Running tests for ${TEST_NAME}:"
    cd ${SUBDIR}
    go test -coverprofile=${COVERAGE_RESULTS}

    echo -e "\n>>> test coverage:"
    go tool cover -func ${COVERAGE_RESULTS}

    echo -e "\nSaved results as test_coverage/${TEST_NAME}.html"
    go tool cover -html ${COVERAGE_RESULTS} -o ${COVERAGE_HTML}
done
