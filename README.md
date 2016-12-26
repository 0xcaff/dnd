DnD
===

A file transfer tool.

Demo
----

![Demo][demo]


Usage
-----

Get

    go get github.com/caffinatedmonkey/dnd

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

[demo]: https://caffinatedmonkey.github.io/dnd/demo.gif
[latest]: https://github.com/caffinatedmonkey/dnd/releases
