# icsgo

[![Build Status](https://travis-ci.org/freechessclub/icsgo.svg)](https://travis-ci.org/freechessclub/icsgo)
[![GoDoc](https://godoc.org/github.com/freechessclub/icsgo?status.svg)](https://godoc.org/github.com/freechessclub/icsgo)
[![GoReportCard](https://goreportcard.com/badge/freechessclub/icsgo)](https://goreportcard.com/report/github.com/freechessclub/icsgo)

icsgo is a Go client library to connect to [Internet Chess Server (ICS)](https://en.wikipedia.org/wiki/Internet_chess_server). An ICS provides support for playing, watching and discussing chess games. Two popular, and among the earliest, examples of ICS include:
* [Free Internet Chess Server (FICS)](http://www.freechess.org/)
* [Internet Chess Club (ICC)](http://www.chessclub.com/)

Although the ICS protocol is a simple variant of the TELNET protocol, it has
not been standardised. As such, this has led to different ICS servers
implementing non-standard extensions to the protocol.

Installation
------------
To install this package, you need to have a working installation of Go.
Assuming you have Go properly installed, simply do:

```
go get -u github.com/freechessclub/icsgo
```

Usage
-----
If you're using Go modules (Go 1.11+), this library can be used by simply
importing `"github.com/freechessclub/icsgo"` in your application. Using the usual Go commands `go [build|run|test]` will automatically download the required dependencies.

Documentation
-------------
* See [godoc](https://godoc.org/github.com/freechessclub/icsgo) for package documentation.
* Examples on how to use the library can be found in the [examples](examples/) directory.

