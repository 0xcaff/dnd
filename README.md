DnD
===

A tool to simply transfer files.

Usage
-----

Get

    go get github.com/caffinatedmonkey/dnd

Run

    dnd

Now navigate to `host:port` and uploaded files will be placed in the current
directory.

Build
-----
If building an stand alone binary, the static assets must be packaged.

    rice embed-go
    go build

