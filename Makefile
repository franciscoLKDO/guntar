TEST_RESULTS_DIR?=`pwd`/test/results
TEST_RESULTS_COVERAGE_REPORT_DIR?=${TEST_RESULTS_DIR}/coverage
TEST_RESULTS_COVERAGE_REPORT_COV?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.cov
TEST_RESULTS_COVERAGE_REPORT_HTML?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.html
TEST_TIMEOUT?=100s
TEST_REPEAT_COUNT?=3
APP_VERSION?=$(shell go run . version)-dev

GENERATE_ID=$(shell docker create guntar)
SET_DOCKER_ID = $(eval DOCKER_ID=$(GENERATE_ID))


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

build:
	docker build . -t guntar --build-arg APP_VERSION=${APP_VERSION}

vhs-tape: build vhs/*.tape
	$(SET_DOCKER_ID)
	docker cp ${DOCKER_ID}:/guntar ./vhs && docker rm ${DOCKER_ID}
	for file in $^ ; do \
		docker run --rm \
			-u `id -u`:`id -u` \
			-v `pwd`/test/mytarfolder.tar:/mytarfolder.tar \
			-v `pwd`/vhs:/vhs \
			-e TARFILE=/mytarfolder.tar -e OUTDIR=./extracted \
			ghcr.io/charmbracelet/vhs $${file##*/} ; \
	done
	rm ./vhs/guntar
	rm -r ./vhs/extracted


install-binary: install
	go install