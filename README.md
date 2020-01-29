# goms

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/goms)](https://travis.mpi-internal.com/Yapo/goms)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/goms)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/goms?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/goms)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/goms?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/goms)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/goms)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/goms)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/goms)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/goms)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

Goms is the official golang microservice template for Yapo.

## A few rules

* Goms was built following [Clean Architecture](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) so, please, familiarize yourself with it and let's code great code!

* Goms has great [test coverage](https://quality-gate.mpi-internal.com/#/Yapo/goms) and [examples](https://github.mpi-internal.com/Yapo/goms/search?l=Go&q=func+Test&type=&utf8=%E2%9C%93) of how good testing can be done. Please honor the effort and keep your test quality in the top tier.

* Goms is not a silver bullet. If your service clearly doesn't fit in this template, let's have a [conversation](mailto:dev@schibsted.cl)

* [README.md](README.md) is the entrypoint for new users of your service. Keep it up to date and get others to proof-read it.

## How to run the service

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/goms.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd goms
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
  

## Creating a new service

* Create a repo for your new service on: https://github.mpi-internal.com/Yapo
* Rename your goms dir to your service name:
  ```
  $ mv goms YourService
  ```
* Update origin: 
  ```
  # https://help.github.com/articles/changing-a-remote-s-url/
  $ git remote set-url origin git@github.mpi-internal.com:Yapo/YourService.git
  ```

* Replace every goms reference to your service's name:
  ```
  $ git grep -l goms | xargs sed -i.bak 's/goms/yourservice/g'
  $ find . -name "*.bak" | xargs rm
  ```

* Go through the code examples and implement your service
  ```
  $ git grep -il fibonacci
  README.md
  cmd/goms/main.go
  pkg/domain/fibonacci.go
  pkg/domain/fibonacci_test.go
  pkg/interfaces/handlers/fibonacci.go
  pkg/interfaces/handlers/fibonacci_test.go
  pkg/interfaces/loggers/fibonacciInteractorLogger.go
  pkg/interfaces/repository/fibonacci.go
  pkg/interfaces/repository/fibonacci_test.go
  pkg/usecases/getNthFibonacci.go
  pkg/usecases/getNthFibonacci_test.go
  ```

* Enable TravisCI
  - Go to your service's github settings -> Hooks & Services -> Add Service -> Travis CI
  - Fill in the form with the credentials you obtain from https://travis.mpi-internal.com/profile/
  - Sync your repos and organizations on Travis
  - Make a push on your service
  - The push should trigger a build. If it didn't ensure that it is enabled on the travis service list
  - Enjoy! This should automatically enable quality-gate reports and a few other goodies

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

### GET  /fibonacci
Implements the Fibonacci Numbers with Clean Architecture

#### Request
{
	"n": int - Ask for the nth fibonacci number
}

#### Response

```javascript
200 OK
{
	"Result": int - The nth fibonacci number
}
```

#### Error response
```javascript
400 Bad Request
{
	"ErrorMessage": string - Explaining what went wrong
}
```

### GET  /user/basic-data?mail=[user_mail]
Returns the essential user data. It is in communication with the Profile Microservice. The main goal of this endpoint is to be used for a basic Pact Test.

#### Request

No additional parameters

#### Response

```javascript
200 OK
{
    "fullname": Full name of the user,
    "cellphone": The user´s cellphone,
    "gender": The user gender,
    "country": The country where the user lives (Currently only Chile is Available),
    "region": The region where the user lives,
    "commune": The commune where the user lives,
}
### Contact
dev@schibsted.cl

## Kubernetes

Kubernetes and Helm have to be installed in your machine.
If you haven't done it yet, you need to create a secret to reach Artifactory.
`kubectl create secret docker-registry containers-mpi-internal-com -n <namespace> --docker-server=containers.mpi-internal.com --docker-username=<okta_username> --docker-password=<artifactory_api_key> --docker-email=<your_email>`

### Helm Charts

1. You need to fill out the ENV variables in the k8s/goms/templates/configmap.yaml file.
2. You should fill out the *tag*, and *host* under hosts to your namespace.
3. Add this host name to your /etc/hosts file with the correct IP address (127.21.5.11)
4. You run `helm install -n <name_of_your_release> k8s/goms`
5. Check your pod is running with `kubectl get pods`
6. If you want to check your request log `kubectl logs <name_of_your_pod>`
