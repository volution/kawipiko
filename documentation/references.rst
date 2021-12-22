

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




References
==========


.. [CDB]
    * `CDB <https://en.wikipedia.org/wiki/Cdb_(software)>`__ (@Wikipedia);
    * `cdb <http://cr.yp.to/cdb.html>`__ (project website, reference implementation by DJB);
    * `cdb <https://github.com/colinmarc/cdb>`__ (project @GitHub, pure Go implementation, used by ``kawipiko`` with patches;)
    * `Constant Database Internals <http://www.unixuser.org/~euske/doc/cdbinternals/index.html>`__ (article);
    * `Benchmarking LevelDB vs. RocksDB vs. HyperLevelDB vs. LMDB Performance for InfluxDB <https://www.influxdata.com/blog/benchmarking-leveldb-vs-rocksdb-vs-hyperleveldb-vs-lmdb-performance-for-influxdb/>`__ (article);
    * `Badger vs LMDB vs BoltDB: Benchmarking key-value databases in Go <https://blog.dgraph.io/post/badger-lmdb-boltdb/>`__ (article);
    * `Benchmarking BDB, CDB and Tokyo Cabinet on large datasets <https://www.dmo.ca/blog/benchmarking-hash-databases-on-large-data/>`__ (article);
    * `TinyCDB <http://www.corpit.ru/mjt/tinycdb.html>`__ (fork project);
    * `tinydns <https://cr.yp.to/djbdns/tinydns.html>`__ (DNS server using CDB);
    * `qmail <https://cr.yp.to/qmail.html>`__ (SMTP server using CDB);


.. [Go]
    * `Go <https://en.wikipedia.org/wiki/Go_(programming_language)>`__ (@Wikipedia);
    * `Go <https://golang.com/>`__ (project website);


.. [fasthttp]
    * `fasthttp <https://github.com/valyala/fasthttp>`__ (project @GitHub);
    * high performance HTTP server implementation;  (alternative to Go's ``net/http`` implementation;)
    * supports HTTP/1 (with or without TLS);
    * used by ``kawipiko``;


.. [quic-go]
    * `quic-go <https://github.com/lucas-clemente/quic-go>`__ (project @GitHub);
    * supports HTTP/3 (over QUIC);
    * used by ``kawipiko``;


.. [Zopfli]
    * `Zopfli <https://en.wikipedia.org/wiki/Zopfli>`__ (@Wikipedia);
    * `Zopfli <https://github.com/google/zopfli>`__ (project @GitHub, reference implementation by Google);
    * `Zopfli <https://github.com/foobaz/go-zopfli>`__ (project @GitHub, pure Go implementation, used by ``kawipiko``);


.. [Brotli]
    * `Brotli <https://en.wikipedia.org/wiki/Brotli>`__ (@Wikipedia);
    * `Brotli <https://github.com/google/brotli>`__ (project @GitHub, reference implementation by Google);
    * `Brotli <https://github.com/andybalholm/brotli>`__ (project @GitHub, pure Go implementation, used by ``kawipiko``);
    * `Results of experimenting with Brotli for dynamic web content <https://blog.cloudflare.com/results-experimenting-brotli/>`__ (article);


.. [Blake3]
    * `Blake3 <https://en.wikipedia.org/wiki/BLAKE_(hash_function)>`__ (@Wikipedia);
    * `Blake3 <https://github.com/BLAKE3-team/BLAKE3>`__ (project @GitHub, reference implementation);
    * `Blake3 <https://github.com/zeebo/blake3>`__ (project @GitHub, pure Go implementation, used by ``kawipiko``);


.. [Bolt]
    * `bolt <https://github.com/boltdb/bolt>`__ (project @GitHub, original pure Go implementation);
    * `bbolt <https://github.com/etcd-io/bbolt>`__ (project @GitHub, forked pure Go implementation, used by ``kawipiko``);


.. [wrk]
    * `wrk <https://github.com/wg/wrk>`__ (project @GitHub);
    * modern HTTP benchmarking tool;
    * multi threaded, implemented in C, with event loop and Lua support;
    * supports HTTP/1 (with and without TLS);


.. [h2load]
    * part of the ``nghttp2`` project;
    * `nghttp2 <https://github.com/nghttp2/nghttp2>`__ (project @GitHub);
    * modern HTTP benchmarking tool;
    * multi threaded, implemented in C, with event loop;
    * supports HTTP/1 (with TLS), HTTP/3, and HTTP/3 (over QUIC);


.. [Netlify]
    * `Netlify <https://www.netlify.com/>`__ (cloud provider);


.. [HAProxy]
    * `HAProxy <https://en.wikipedia.org/wiki/HAProxy>`__ (@Wikipedia);
    * `HAProxy <https://www.haproxy.org/>`__ (project website);
    * reliable high performance TCP/HTTP load-balancer;
    * multi threaded, implemented in C, with event loop and Lua support;


.. [NGinx]
    * `NGinx <https://en.wikipedia.org/wiki/Nginx>`__ (@Wikipedia);
    * `NGinx <https://nginx.org/>`__ (project website);
    * reliable high performance HTTP server;
    * multi threaded, implemented in C, with event loop;


.. [darkhttpd]
    * `darkhttpd <https://unix4lyfe.org/darkhttpd/>`__ (project website);
    * simple static HTTP server;
    * single threaded, implemented in C, with event loop and ``sendfile`` support;


.. [mmap]
    * `Memory mapping <https://en.wikipedia.org/wiki/Memory-mapped_file>`__ (@Wikipedia);
    * `mmap(2) <http://man7.org/linux/man-pages/man2/mmap.2.html>`__ (Linux man page);

