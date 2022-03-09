#!/bin/sh

export GO111MODULE=on
go get github.com/bhojpur/web/@latest
webctl new hello
cd hello
webctl run