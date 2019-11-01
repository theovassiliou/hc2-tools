# hc2-tools

hc2-tools, implemented in go (golang) provides tools to, upload, download and interact with scene on Fibaro HC2 system.

The design goal of the tools where to be easily integrable into IDEs like Visual Studio Code (VSC) or others, to enable comfortable Fibaro Lua script development, testing and deployment.
While there can be very different applications on how to use the tools, [USAGE.md](USAGE.md) has examples on how the development workflow can be enhanced by using the hc2-tools.

It is in public beta; you are free to use it and improve it (do pull requests).

WARNING: THIS SOFTWARE CAN'T BE ERROR FREE, SO USE IT AT YOUR OWN RISK. DON'T USE IT IF YOU HAVEN'T MADE AN ACTUAL BACKUP COPY OF YOUR FIBARO HC2 SYSTEM. IF YOU DO NOT HOW TO DO THIS, PLEASE RECONSIDER TO USE THIS SOFTWARE ANYWAY. I HAVE DONE MY BEST TO MAKE SURE THAT THE TOOLS BEHAVE AS EXPECTED. BUT AGAIN ... USE IT AT YOUR OWN RISK. I AM NOT GIVING ANY KIND OF WARRANTY, NEITHER EXPLICITELY NOR IMPLICITELY.

## Installation binaries

You can download the binaries directly from the [releases](https://github.com/theovassiliou/hc2-tools/releases) section.  Unzip/untar the downloaded archive and copy the files to a location of your choice, e.g. `/usr/local/bin/` on *NIX or MacOS. If you install only the binaries, make sure that they are accessible from the command line. Ideally, they are accessible via `$PATH` or `%PATH%`, respectively.

### Configuring your installation

hc2-tools need access to your Fibaro HC2 system. You can configure and test your installation just by

`hc2DownloadScene -u <fibaroHc2Login> -p <secretPassword> --url http://<ip.address.of.hc2> -i -t`

This tests the connection to your Fibaro HC2 system and creates a config-file in `~/.hc2-tool/*` so that you do not have to reenter the information in subsequent to the hc2-tools calls.
Of course, you have to replace `<fibaroHc2Login>` and `<secretPassword>` with your credentials, and `<ip.address.of.hc2>` with the IP or DNS of your Fibaro HC2 system.

To test the successfull configuration just use

`hc2DownloadScene -t` which should get you the same result as above. All h2-tools use the same configuration file, so you don't have to configure the individually. Actually you could perfom this configuration steps, with any of the tools.

If you would like to learn about more about the technology in the background take a look at [TECHNOLOGY.md](TECHNOLOGY.md).

## Installation From Source

hc2-tools requires golang version 1.3 or newer, the Makefile requires GNU make.

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

There is no particular requirement beyong the fact that you should have a working go installation.

[Install Go](https://golang.org/doc/install) >=1.13

### Installing

Download hc2-tools source by running

```shell
go get -u github.com/theovassiliou/hc2-tools
```

This gets you your copy of hc2-tools installed under
`$GOPATH/src/github.com/theovassiliou/hc2-tools`

Run `make` from the source directory by running

```shell
cd $GOPATH/src/github.com/theovassiliou/hc2-tools
make all
```

to compile and build the three executables

* hc2DownloadScene - [README](cmd/hc2DownloadScene/README.md)
* hc2UploadScene - [README](cmd/hc2UploadScene/README.md)
* hc2SceneInteraction - [README](cmd/hc2SceneInteraction/README.md)

and run

```shell
make go-install
```

to install a copy of the executables into `$GOPATH/bin`

To test whether build integrates with your Fibaro system use

```shell
hc2DownloadScene -u <fibaroHc2Login> -p <secretPassword> --url http://<ip.address.of.hc2> -i -t

Successfull connected to ...
  Name         : Hal
  Serial       : HC2-033533
  IP           : <ip.address.of.hc2>
  Version      : 4.560
  ZWaveVersion : 3.67

and logged in as:
  User:         <fibaroHc2Login>
  Type:         superuser
```

This tests the connection to your Fibaro HC2 system and `-i` creates a config-file in `~/.hc2-tool/*` so that you do not have to reenter the information in subsequent to the hc2-tools calls.
Of course you have to replace `<fibaroHc2Login>` and `<secretPassword>` with your credentials, and `<ip.address.of.hc2>` with the IP or DNS of your Fibaro HC2 system.

To test the successfull configuration just use

`hc2DownloadScene -t` which should get you the same result as above.

Now make your head-start by calling

```shell
hc2DownloadScene
```

which downloads all scripts on the Fibaro HC2 into the default created directory `./download/`

## Running the tests

We are using to different make targets for running tests.

```shell
make test

go test -short ./...
?    github.com/theovassiliou/hc2-tools/cmd/expandRequire [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2DownloadScene [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2SceneInteract [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2UploadScene [no test files]
ok   github.com/theovassiliou/hc2-tools/pkg 0.036s
```

executes all short package tests, while

```shell
make test-all
go vet $(go list ./...)
go test ./...
?    github.com/theovassiliou/hc2-tools/cmd/expandRequire [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2DownloadScene [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2SceneInteract [no test files]
?    github.com/theovassiliou/hc2-tools/cmd/hc2UploadScene [no test files]
ok   github.com/theovassiliou/hc2-tools/pkg 0.036s
```

executes in addition `go vet`on the package. Before committing to the code base please use `make test-all` to ensure that all tests pass.

### Break down into end to end tests

After creating your configuration call `hc2DownloadScene` without `-u -p` parameters.

```shell
hc2DownloadScene -t

Successful connected to ...
  Name         : Hal
  Serial       : HC2-033533
  IP           : 192.10.66.55
  Version      : 4.560
  ZWaveVersion : 3.67

and logged in as:
  User:         specialuser@mydomain.com
  Type:         superuser
```

to test whether command can execute correctly.

## Deployment

After running

```shell
make install

go build -ldflags " -X main.commit=99f909d -X main.branch=master" ./cmd/hc2UploadScene
go build -ldflags " -X main.commit=99f909d -X main.branch=master" ./cmd/hc2DownloadScene
go build -ldflags " -X main.commit=99f909d -X main.branch=master" ./cmd/hc2SceneInteract
mkdir -p /usr/local/bin/
cp hc2UploadScene /usr/local/bin/
cp hc2DownloadScene /usr/local/bin/
cp hc2SceneInteract /usr/local/bin/
```

you can find your executables in `/usr/local/bin`. Make sure `/usr/local/bin/` is in your path.

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/theovassiliou/hc2-tools/tags).

## Authors

* **Theo Vassiliou** - *Initial work* - [Theo Vassiliou](https://github.com/theovassiliou)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

Thanks to all the people out there that produce amazing open-source software, which supported the creation of this piece of software. In particular I wasn't only able to use libraries etc. But also, to learn and understand golang better. In particular I wanted to thank

* [Jaime Pillora](https://github.com/jpillora) for [jpillora/opts](https://github.com/jpillora/opts). Nice piece of work!
* [InfluxData Team](https://github.com/influxdata) for [influxdata/telegraf](https://github.com/influxdata/telegraf). Here I learned a lot for Makefile writing and release building in particular.
* Inspiration and motivation to develop this tool I got from the [ZeroBrane](https://studio.zerobrane.com/) Lua Development Environment.
* [PurpleBooth](https://gist.github.com/PurpleBooth) for the well motivated [README-template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)

***

## History

This project has been developed as I was seeking for a way to upload scenes from my favorite development Lua development environment to the Fibaro HC2 system. Finally I came up with the idea to upload a scene whenever I do a `git commit`. For this I needed a cmd line tool can be integrated as `commithook` into the git repository.

With this I could solve two problems at a single time.

1. Enforcing a version control system, e.g. git
2. Automatically uploading the modified script

After implementing a first version, new ideas emerged, so for example retrieving debug messages where implemented.
