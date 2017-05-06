# Go Commons

Find yourself needing the same packages over and over?

Me too.

Instead of re-writting them for each project/client/company I decided it was time to do it one more time in public.

## Getting Started

To use these packages simply:
```
go get github.com/corsc/go-commons/
```
(vendoring recommended)

To contribute:

* Fork
* PR
* Enjoy!

### Prerequisites

* Go 1.8
* (optional) [GoMetaLinter](https://github.com/alecthomas/gometalinter)
* (optional) [My GoMetaLinter Config](https://raw.githubusercontent.com/corsc/PersonalTools/master/go-scripts/gometa-config.json)

## Running the tests

Nothing special, standard `go test ./...` will get the job down.

If me, you want the fastest possible tests, I would skip vendor by using:
```
go test $(go list ./... | grep -v /vendor)
```

### Lint checking contributions

Please check all PRs before sending them using the following settings

```
gometalinter --config=gometa-config.json ./...
```

## Authors

* **Corey Scott** - [corsc](https://github.com/corsc)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

