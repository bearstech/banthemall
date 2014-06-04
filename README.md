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

Use it
------

Find a `GeoIP.dat` index from maxmind.

    tail -f /var/log/apache2/access/log | ./banthemall

More informations are shown every 30s :

 - Country code
 - IP
 - Number of request per 30s
 - [Spamhaus status](http://www.spamhaus.org/zen/). - is good, PBL is manageable, other are drama.


License
-------

3 terms BSD licence, © Mathieu Lecarme 2014.
