APP_NAME:=guntar
TEST_RESULTS_DIR?=`pwd`/test/results
TEST_RESULTS_COVERAGE_REPORT_DIR?=${TEST_RESULTS_DIR}/coverage
TEST_RESULTS_COVERAGE_REPORT_COV?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.cov
TEST_RESULTS_COVERAGE_REPORT_HTML?=${TEST_RESULTS_COVERAGE_REPORT_DIR}/coverage.html
TEST_TIMEOUT?=100s
TEST_REPEAT_COUNT?=3
APP_VERSION?=$(shell go run . version)


DOCKER_TARGET?=prod
DOCKER_STAGE?=dev
DOCKER_TAG?=${APP_VERSION}-${DOCKER_STAGE}
DOCKER_RESULTS_DIR=/app/test/results
COMMIT_ID?=""


test-setup:
	mkdir -p ${TEST_RESULTS_DIR}
	mkdir -p ${TEST_RESULTS_COVERAGE_REPORT_DIR}

test: test-setup
	@go test -timeout ${TEST_TIMEOUT} -v -count ${TEST_REPEAT_COUNT} ./... -coverpkg=./... -coverprofile=${TEST_RESULTS_COVERAGE_REPORT_COV}
	@go tool cover -html=${TEST_RESULTS_COVERAGE_REPORT_COV} -o ${TEST_RESULTS_COVERAGE_REPORT_HTML}
	@echo "Coverage report available here: ${TEST_RESULTS_COVERAGE_REPORT_HTML}\n"

install-dependencies:
	go mod tidy
	go mod download && go mod verify

install-binary: install
	go install

build-image:
	docker build . -t ${APP_NAME}:${DOCKER_TAG} --build-arg APP_VERSION=${APP_VERSION} COMMIT_ID=${COMMIT_ID} --target ${DOCKER_TARGET}

test-ci:
	docker run -t -v ${TEST_RESULTS_DIR}:${DOCKER_RESULTS_DIR} ${APP_NAME}:${DOCKER_TAG} make test

# Handle gifs build.
GIF_DIR:=vhs
TARBALL:=/archive_gif.tar
OUTDIR:=extracted

OUTPUTS_GIF = $(patsubst ${GIF_DIR}/%.tape,${GIF_DIR}/%.gif,$(wildcard ${GIF_DIR}/*.tape))

# This is building gifs inside vhs container, do not use it directly, you must use build-% target instead
# env variables are handled here and can be used in tapes.
# Guntar docker image must be present on host.
${GIF_DIR}/%.gif: ${GIF_DIR}/%.tape
	@echo "Run gif creation for $@ from tape $(^F) \n"
	@docker run --rm \
		-u `id -u`:`id -u` \
		-v `pwd`/test/mytarfolder.tar:${TARBALL} \
		-v `pwd`/${GIF_DIR}:/vhs \
		-e TARBALL=${TARBALL} -e OUTDIR=${OUTDIR} \
		ghcr.io/charmbracelet/vhs $(^F)
	@echo "New gif here: $@\n"

pre-build-gif:
	@echo "Build app\n"
	@$(eval DOCKER_ID=$(shell docker create ${APP_NAME}))
	@docker cp ${DOCKER_ID}:/${APP_NAME} ${GIF_DIR} && docker rm ${DOCKER_ID}

post-build-gif:
	@echo "Clean files\n"
	@rm -rf ${GIF_DIR}/${OUTDIR}_*;
	@rm -f ./${GIF_DIR}/${APP_NAME}

# Build all updated gifs
build-gifs: pre-build-gif $(OUTPUTS_GIF) 
	@$(MAKE) post-build-gif

# Build specific gif if updated
build-gif-%: pre-build-gif ${GIF_DIR}/%.gif
	@$(MAKE) post-build-gif

# Remove all existing gifs
remove-gifs:
	rm ${GIF_DIR}/*.gif

.PHONY: pre-build-gif post-build-gif build-all build-% remove-gif