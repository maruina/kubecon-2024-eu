#!/bin/bash

TEST_FILE=$(mktemp)
SERVICE="kubecon-2024-eu"

# Running gotestsum to execute a compiled test binary (created with go test -c)
# See https://github.com/gotestyourself/gotestsum?tab=readme-ov-file#executing-a-compiled-test-binary
# See https://github.com/gotestyourself/gotestsum/blob/main/.project/docs/running-without-go.md
echo "INFO: running E2E tests..."
exit_code=0
if ! ./gotestsum --format standard-verbose --junitfile "${TEST_FILE}" --raw-command -- ./test2json -t -p e2e ./e2e.test "$@"; then
    exit_code=1
fi

if [[ "${DATADOG_API_KEY}" == "PLACEHOLDER" ]]; then
  echo "INFO: DATADOG_API_KEY not set, skip sending tests results"
  exit "${exit_code}"
fi

# Collecting git metadata
# See https://docs.datadoghq.com/tests/setup/junit_xml/?tab=macos#collecting-git-metadata
# DD_GIT_TAG and DD_GIT_COMMIT_SHA are created at build time in the Dockerfile
export DD_GIT_REPOSITORY_URL="https://github.com/maruina/kubecon-2024-eu"
export DD_GIT_BRANCH="main"
export DATADOG_SITE="datadoghq.eu"

# Upload test results to Datadog even on test failure
if ! datadog-ci junit upload \
    --service "${SERVICE}" \
    --env "${ENV:-undefined}" \
    --tags cluster:"${CLUSTER_NAME:-undefined}" \
    "${TEST_FILE}"; then
    exit_code=1
fi

exit "${exit_code}"
