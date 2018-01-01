DnD
===

[![Build Status][build-status-image]][build-status]

A file transfer tool.

Demo
----

![Demo][demo]


Usage
-----

Get

    go get github.com/0xcaff/dnd

[or download the latest release (you don't need go to run it)][latest]

Run

    dnd

Now navigate to `host:port`, select or drop files onto the page. The files will
be uploaded to the current directory.

Build
-----
If building an stand alone binary for distribution, the static assets must be
packaged.

    go generate
    go build

[demo]: https://0xcaff.github.io/dnd/demo-progress.gif
[latest]: https://github.com/0xcaff/dnd/releases

[build-status-image]: https://travis-ci.org/0xcaff/dnd.svg?branch=master
[build-status]: https://travis-ci.org/0xcaff/dnd
