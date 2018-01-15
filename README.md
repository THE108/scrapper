# ScrappeR

Testing automation tool

Usage
=====

```
./scrapper -h
Usage of ./scrapper:
  -file string
    	Config file name (yaml). Use to validate shipping info according to config.
  -url string
    	Url to parse. Use to get shipping info from PDP page.
```

There are two things this tool could be used for:
  - Validation: just provide config file name (see `example_config.yaml`)
  
  - Info retrieval: just provide an URL to parse
  
Examples
========

Validation:
```
$ ./scrapper -file config.yaml
```

Info retrieval:
```
$ ./scrapper -url http://test.com
```

Installation
============

 - Install `golang`
 
See https://golang.org/doc/install

 - Install `dep` (https://github.com/golang/dep)

```
$ brew install dep
$ brew upgrade dep
```

or
```
$ go get -u github.com/golang/dep/cmd/dep
```

 - Download dependencies:
 
Just run:
```
$ dep ensure
```

 - Compile `scrapper`
 
```
$ go build scrapper
```

Licence
=======

MIT