# device-controller

`device-controller` is a component within the Nalej platform that resides in the application cluster and is in charge of having updated information about the devices running in a cluster. 

## Getting Started

The `device-controller` component has three main functionalities:
* Ping
* Ping registration
* Latency check

### Prerequisites

* [cluster-api](https://github.com/nalej/cluster-api)
* [login-api](https://github.com/nalej/login-api)

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

No test files are available for this repository at this moment.

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/device-controller/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/device-controller/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
