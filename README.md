# premium-carousel-api

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/premium-carousel-api)](https://travis.mpi-internal.com/Yapo/premium-carousel-api)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/premium-carousel-api)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/premium-carousel-api?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/premium-carousel-api)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/premium-carousel-api?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/premium-carousel-api)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/premium-carousel-api)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/premium-carousel-api)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/premium-carousel-api)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/premium-carousel-api)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

premium-carousel-api needs a description here.


## How to run premium-carousel-api

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/premium-carousel-api.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd premium-carousel-api
  $ make start
  ```

* To get a list of available commands:

  ```
  $ make help
  Targets:
    test                 Run tests and generate quality reports
    cover                Run tests and output coverage reports
    coverhtml            Run tests and open report on default web browser
    checkstyle           Run gometalinter and output report as text
    setup                Install golang system level dependencies
    build                Compile the code
    run                  Execute the service
    start                Compile and start the service
    fix-format           Run gofmt to reindent source
    info                 Display basic service info
    docker-build         Create docker image based on docker/dockerfile
    docker-publish       Push docker image to containers.mpi-internal.com
    docker-attach        Attach to this service's currently running docker container output stream
    docker-compose-up    Start all required docker containers for this service
    docker-compose-down  Stop all running docker containers for this service
    help                 This help message
  ```

* If you change the code:

  ```
  $ make start
  ```

* How to run the tests

  ```
  $ make [cover|coverhtml]
  ```

* How to check format

  ```
  $ make checkstyle
  ```

## Endpoints
### GET  /healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

## Contact
dev@schibsted.cl

## Kubernetes

Kubernetes and Helm have to be installed in your machine.
If you haven't done it yet, you need to create a secret to reach Artifactory.
`kubectl create secret docker-registry containers-mpi-internal-com -n <namespace> --docker-server=containers.mpi-internal.com --docker-username=<okta_username> --docker-password=<artifactory_api_key> --docker-email=<your_email>`

### Helm Charts

1. You need to fill out the ENV variables in the k8s/premium-carousel-api/templates/configmap.yaml file.
2. You should fill out the *tag*, and *host* under hosts to your namespace.
3. Add this host name to your /etc/hosts file with the correct IP address (127.21.5.11)
4. You run `helm install -n <name_of_your_release> k8s/premium-carousel-api`
5. Check your pod is running with `kubectl get pods`
6. If you want to check your request log `kubectl logs <name_of_your_pod>`
