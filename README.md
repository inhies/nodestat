Nodestat
========

Nodestat is a tool that allows you to easily gather statistics about your cjdns
node. It currently presents a very simple HTTP JSON API for getting information. 

**Nodestat is currently undergoing heavy development and as such the structure,
format, and content of the data it provides is subject to change at any time.**


Installation
------------

To install nodestat, run `go get github.com/inhies/nodestat`. You will need to
create a configuration file. Simply copy the sample file to `nodestat.conf` and
place it in the same directory as your nodestat executable. 


Usage
-----

By default, nodestat will present three endpoints: 

* `/node` Gives you information regarding cjdns itself
* `/peers` Gives you information regarding your connected peers
* `/all` Combines node and peers output
