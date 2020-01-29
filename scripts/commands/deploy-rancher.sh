#!/usr/bin/env bash

# Only tags of the form <version@target> will be deployed
if [[ ! "${TRAVIS_BRANCH}" =~ "@" ]]; then
	echo "This tag is not a deploy tag [${TRAVIS_BRANCH}]"
	exit 0
fi

# In case we are in travis, we will use cached docker environment.
if [[ -n "$TRAVIS" ]]; then
    DOCKER_COMMAND=container_cache
else
    DOCKER_COMMAND=docker
fi

echo Preparing environment
GIT_COMMIT=$(git rev-list -n 1 ${TRAVIS_BRANCH})
TAG_NAME=${TRAVIS_BRANCH%%@[[:alnum:]]*}
TAG_DEPLOY_ENV=${TRAVIS_BRANCH##[[:alnum:]]*@}
DOCKER_RUN_ENV="                         \
  -e TAG_NAME=${TRAVIS_BRANCH}           \
  -e TARGET=${TAG_DEPLOY_ENV}            \
  -e IMAGE=${DOCKER_IMAGE}:${TAG_NAME}   \
  -e SERVICE_REPO=Yapo/${APPNAME}        \
  -e ENV_REPO=${RANCHER_ENV_REPO}        \
  -e APPNAME                             \
  -e BUILD_ID=${TRAVIS_BUILD_ID}         \
  -e BUILD_NUMBER=${TRAVIS_BUILD_NUMBER} \
  -e REPO_SLUG=${TRAVIS_REPO_SLUG}       \
  -e GIT_COMMIT                          \
  -e GITHUB_ACCESS_TOKEN                 \
  -e RANCHER_DEV_API_KEY                 \
  -e RANCHER_DEV_SECRET_KEY              \
  -e RANCHER_PRE_API_KEY                 \
  -e RANCHER_PRE_SECRET_KEY              \
  -e RANCHER_PRO_API_KEY                 \
  -e RANCHER_PRO_SECRET_KEY              \
  -e DEVHOSE_KEY                         \
  -e ARTIFACTORY_USER                    \
  -e ARTIFACTORY_PWD                     \
  -e ARTIFACTORY_CONTEXT                 "

echo Starting deploy process
${DOCKER_COMMAND} run ${DOCKER_RUN_ENV} -ti ${RANCHER_DEPLOY_IMAGE}
