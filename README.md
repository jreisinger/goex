Introduction to Go Programming
==============================

Code from my [Introduction to Go
Programming](http://shop.oreilly.com/product/0636920035305.do) course.

sample
======

The sample directory contains the following:

1. x.go, picdumidi.jpg, degraded.jpg - code and sample images used in "The
standard cgo package"

2. theprelude, wordsworth - poetry samples

3. src/hello - blank starter program

4. src/poetry - The poetry package and tests

5. src/shuffler - The shuffler package

6. config - The JSON config file

Usage
=====

```
cd sample
export GOPATH=`pwd`

go test poetry

go install hello
./bin/hello

# in different terminal window/tab
curl localhost:8080/poem?name=wordsworth | jq
curl localhost:8080/poem?name=theprelude | jq
```
