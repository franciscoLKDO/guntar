TEST_RESULTS_DIR?=`pwd`/test/results
TEST_RESULTS_COVERAGE_REPORT_DIR?=${TEST_RESULTS_DIR}/coverage
TEST_RESULTS_COVERAGE_REPORT_COV?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.cov
TEST_RESULTS_COVERAGE_REPORT_HTML?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.html
TEST_TIMEOUT?=100s
TEST_REPEAT_COUNT?=10

test-setup:
	mkdir -p ${TEST_RESULTS_DIR}
	mkdir -p ${TEST_RESULTS_COVERAGE_REPORT_DIR}

test: test-setup
	@go test -timeout ${TEST_TIMEOUT} -v -count ${TEST_REPEAT_COUNT} ./... -coverpkg=./... -coverprofile=${TEST_RESULTS_COVERAGE_REPORT_COV}
	@go tool cover -html=${TEST_RESULTS_COVERAGE_REPORT_COV} -o ${TEST_RESULTS_COVERAGE_REPORT_HTML}
	@printf "Coverage report available here: ${TEST_RESULTS_COVERAGE_REPORT_HTML}\n"

install:
	go mod tidy
	go mod download && go mod verify

install-binary: install
	go install