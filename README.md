Nodestat
========

Nodestat is a tool that allows you to easily gather statistics about your cjdns
node. It currently presents a very simple HTTP JSON API for getting information. 

**Nodestat is currently undergoing heavy development and as such the structure,
format, and content of the data it provides is subject to change at any time.**

Installation
------------

Eventually it will be as easy as `go get github.com/inhies/nodestat` but until
Emery and I are done refactoring the go-cjdns package, some trickery is needed.
You will need to copy the go-cjdns repository in your GOPATH to a new folder
named go-cjdns-refactor, then checkout the `refactor` branch. Here is how I do
it:

    cd $GOPATH/src/github.com/inhies
    cp -R go-cjdns go-cjdns-refactor
    cd go-cjdns-refactor
    git checkout -b refactor
    git pull origin refactor

You now have a copy of the go-cjdns package using the newly refactored code
available as `go-cjdns-refactor`. You should now be able to use `go get` to
fetch nodestat. If this doesn't work, sorry. Figure it out.

Once you have nodestat, copy the sample config to a new file named
`nodestat.conf` and then use `go build` to build it. `./nodestat` to run it. 


Usage
-----

By default, nodestat will present three endpoints: 

* `/node` Gives you information regarding cjdns itself
* `/peers` Gives you information regarding your connected peers
* `/all` Combines node and peers output
