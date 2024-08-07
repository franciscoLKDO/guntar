TEST_RESULTS_DIR?=`pwd`/test/results
TEST_RESULTS_COVERAGE_REPORT_DIR?=${TEST_RESULTS_DIR}/coverage
TEST_RESULTS_COVERAGE_REPORT_COV?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.cov
TEST_RESULTS_COVERAGE_REPORT_HTML?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.html
TEST_TIMEOUT?=100s
TEST_REPEAT_COUNT?=3
APP_VERSION?=$(shell go run . version)-dev
APP_NAME:=guntar


test-setup:
	mkdir -p ${TEST_RESULTS_DIR}
	mkdir -p ${TEST_RESULTS_COVERAGE_REPORT_DIR}

test: test-setup
	@go test -timeout ${TEST_TIMEOUT} -v -count ${TEST_REPEAT_COUNT} ./... -coverpkg=./... -coverprofile=${TEST_RESULTS_COVERAGE_REPORT_COV}
	@go tool cover -html=${TEST_RESULTS_COVERAGE_REPORT_COV} -o ${TEST_RESULTS_COVERAGE_REPORT_HTML}
	@printf "Coverage report available here: ${TEST_RESULTS_COVERAGE_REPORT_HTML}\n"

install-dependencies:
	go mod tidy
	go mod download && go mod verify

install-binary: install
	go install

build-image:
	docker build . -t ${APP_NAME} --build-arg APP_VERSION=${APP_VERSION}


# Handle gif targets, 
GIF_DIR:=vhs
TARBALL:=/mytarfolder.tar
OUTDIR:=./extracted
INPUTS_TAPE = $(wildcard ${GIF_DIR}/*.tape)
OUTPUTS_GIF = $(patsubst ${GIF_DIR}/%.tape,${GIF_DIR}/%.gif,$(INPUTS_TAPE))

all-gif: $(OUTPUTS_GIF)

# TODO find a way to run pre and post rules once  
${GIF_DIR}/%.gif: ${GIF_DIR}/%.tape
	@printf "Build app\n"
	@$(eval DOCKER_ID=$(shell docker create ${APP_NAME}))
	@docker cp ${DOCKER_ID}:/${APP_NAME} ${GIF_DIR} && docker rm ${DOCKER_ID}
	@printf "Run gif creation for $@ from tape $< \n"
	@docker run --rm \
		-u `id -u`:`id -u` \
		-v `pwd`/test/mytarfolder.tar:${TARBALL} \
		-v `pwd`/${GIF_DIR}:/vhs \
		-e TARBALL=${TARBALL} -e OUTDIR=${OUTDIR} \
		ghcr.io/charmbracelet/vhs $(^F)
	@printf "Clean files\n"
	@if [ -d "${GIF_DIR}/extracted" ]; then \
        rm -r ${GIF_DIR}/extracted; \
    fi
	@rm ./${GIF_DIR}/${APP_NAME}
	@printf "Find gif here: $@\n"

clean-gif:
	rm ${GIF_DIR}/*.gif