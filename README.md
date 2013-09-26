Nodestat
========

Nodestat is a tool that allows you to easily gather statistics about your cjdns
node. It currently presents a very simple HTTP JSON API for getting information. 

**Nodestat is currently undergoing heavy development and as such the structure,
format, and content of the data it provides is subject to change at any time.**


Installation
------------

To use nodestat, first clone this repository with `git clone
https://github.com/inhies/nodestat`. Then `cd nodestat`. Copy the example
configuration file to your own working copy with `cp nodestat.conf.sample
nodestat.conf` and build nodestat with `go build`. 

You can start nodestat by running `./nodestat`.


Usage
-----

By default, nodestat will present three endpoints: 

* `/node` Gives you information regarding cjdns itself
* `/peers` Gives you information regarding your connected peers
* `/all` Combines node and peers output
