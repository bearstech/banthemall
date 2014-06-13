Ban them all
============

Realtime log analysis, in a console.

For now, it's a production ready POC.

Build it
--------

Install golang tools.

    go get
    go build

Push the binary on your server.

Cross compilation is very simple with [gxc](https://github.com/robertkrimen/gxc) tool.

Use it
------

Find a `GeoIP.dat` index from maxmind and put in the same folder.

There is some command line options :

    ./banthemall -h

Log are read from STDIN :

    tail -f /var/log/apache2/access/log | ./banthemall

More informations are shown every 10s :

 - Country code
 - IP
 - Number of request per 10s, splitted by http status
 - Number of distinct user agent
 - Number of distinct url
 - [Spamhaus status](http://www.spamhaus.org/zen/). - is good, PBL is manageable, other are drama.


License
-------

3 terms BSD licence, © Mathieu Lecarme 2014.
