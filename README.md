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

To run nodestat, change to it's directory with `cd
$GOPATH/src/github.com/inhies/nodestat`, build it with `go build`, and run it
with `./nodestat`.


Usage
-----

By default, nodestat will present three endpoints: 

* `/node` Gives you information regarding cjdns itself
* `/peers` Gives you information regarding your connected peers
* `/all` Combines node and peers output


Example Output
--------------

The following sections show the output of the `/peers` and `/node` HTTP JSON API
endpoints, respectively.


### Peer Info

```javascript
{
    "gtc8q9d4h4rqvd0wpf7ltt11c85cc5pfffvx6l4rgcdr092yzdp0.k": {
        "PublicKey": "gtc8q9d4h4rqvd0wpf7ltt11c85cc5pfffvx6l4rgcdr092yzdp0.k",
        "State": 2,
        "IsIncoming": false,
        "BytesIn": 612497,
        "BytesOut": 484628,
        "Last": "2013-09-25T12:44:31.679-10:00",
        "SwitchLabel": "0000.0000.0000.0013",
        "IPv6": "fc77:db77:08a3:898d:161f:8de5:66fc:7689",
        "RateIn": 138.88421139965084,
        "RateOut": 119.89995823254979,
        "LastUpdate": "2013-09-25T12:44:31.864414014-10:00"
    },
    "jufvunvsmy08r69dhhw1w170d26qj5kvybdcf45jbkrtw7pb7cl0.k": {
        "PublicKey": "jufvunvsmy08r69dhhw1w170d26qj5kvybdcf45jbkrtw7pb7cl0.k",
        "State": 3,
        "IsIncoming": false,
        "BytesIn": 0,
        "BytesOut": 20664,
        "Last": "2013-09-25T11:28:39.806-10:00",
        "SwitchLabel": "0000.0000.0000.0015",
        "IPv6": "fc19:782e:e609:1c2c:420b:7d2e:cca2:2e3d",
        "RateIn": 0,
        "RateOut": 0,
        "LastUpdate": "2013-09-25T12:44:31.864411162-10:00"
    }
}
```

### Node Info

```javascript
{
    "Memory": 1033176,
    "Angel": {
        "Uptime": "1h15m47s",
        "PercentCPU": 0,
        "PercentMemory": 0
    },
    "Core": {
        "Uptime": "1h15m47s",
        "PercentCPU": 0.7,
        "PercentMemory": 0
    },
    "RateIn": 208.85014166934658,
    "RateOut": 347.73245912192516,
    "BytesIn": 3402765,
    "BytesOut": 2825913
}
```
