
###############
CDB HTTP server
###############


.. contents::
    :depth: 2
    :backlinks: top
    :local:




About
=====

This is a simple static HTTP server written in Go, whose main purpose is to serve (public) static content as efficient as possible.  As such, it basically supports only ``GET`` requests and does not provide features like dynamic content, authentication, reverse proxying, etc.

However it does provide something unique, that no other HTTP server offers:  the static content is served from a CDB_ database with almost zero latency.

CDB_ databases are binary files that provide efficient read-only key-value lookup tables, initially used in some DNS and SMTP servers, mainly for their low overhead lookup operations, zero locking in multi-threaded / multi-process scenarios, and "atomic" multi-record updates.  This also makes them suitable for low-latency static content serving over HTTP, which this project provides.

For a complete list of features please consult the `Features`_ section.

Unfortunately, there are also some tradeoffs as described in the `Limitations`_ section (although none are critical).




Benchmarks
==========


Results
-------

.. note ::

  Bottom line (**even on my 6 years old laptop**):

  * under normal conditions (16 concurrent clients), you get around 36k requests / second, at about 0.5ms latency;
  * under stress conditions (512 concurrent clients), you get arround 32k requests / second, at about 15ms latency;

.. note ::

  Please note that the values under `Thread Stats` are reported per thread.
  Therefore it is best to look at the first two values, i.e. `Requests/sec`.

* 16 connections / 4 threads: ::

    Requests/sec:  36084.51
    Transfer/sec:     16.45MB

      4 threads and 16 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   436.77us  223.21us   3.36ms   81.09%
        Req/Sec     9.07k   499.08    10.27k    72.17%
      Latency Distribution
         50%  390.00us
         75%  481.00us
         90%  669.00us
         99%    1.34ms
      1082680 requests in 30.00s, 493.55MB read

* 512 connections / 4 threads: ::

    Requests/sec:  32773.77
    Transfer/sec:     14.94MB

      4 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    15.84ms   11.04ms  65.68ms   61.64%
        Req/Sec     8.24k     1.76k   15.65k    70.95%
      Latency Distribution
         50%   15.91ms
         75%   23.48ms
         90%   29.63ms
         99%   45.90ms
      986092 requests in 30.09s, 449.52MB read

* 2048 connections / 4 threads: ::

    Requests/sec:  31132.31
    Transfer/sec:     14.19MB

      4 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    98.56ms  163.64ms   4.12s    90.85%
        Req/Sec     7.84k     1.83k   14.43k    68.36%
      Latency Distribution
         50%   57.15ms
         75%   92.95ms
         90%  248.46ms
         99%  671.10ms
      936780 requests in 30.09s, 427.04MB read
      Socket errors: connect 0, read 0, write 1, timeout 0


Notes
-----

The following benchmarks were executed as follows:

* the machine was my personal laptop:  6 years old with an Intel Core i5 2520M (2 cores with 2 threads each), which during the benchmarks (due to a bad fan and dust) it kept entering into thermal throttling;  (i.e. the worst case scenario;)
* the `cdb-http-server` was started with `GOMAXPROCS=4`;  (i.e. 4 threads handling the requests;)
* the `cdb-http-server` was started with `--preload`;  (i.e. the CDB database file was preloaded into memory, thus no disk I/O;)
* the benchmarking tool was wrk_;
* both `cdb-http-server` and `wrk` tools were run on the same machine;
* the benchmark was run over loopback networking (i.e. `127.0.0.1`);
* the served file contains the content ``Hello, World!``;
* the protocol was HTTP  (i.e. no TLS);




Documentation
=============

The project provides two binaries:

* ``cdb-http-server`` -- which serves the static content from the CDB database;
* ``cdb-http-archiver`` -- which creates the CDB database from a source folder holding the static content;




``cdb-http-archiver``
---------------------

::

    >> cdb-http-archiver --help

::

    Usage of cdb-http-archiver:
    --sources <path>
    --archive <path>
    --compress <gzip | brotli | identity>
    --debug




``cdb-http-server``
-------------------

::

    >> cdb-http-server --help

::

    Usage of cdb-http-server:
    --archive <path>
    --archive-inmem      (memory-loaded archive file)
    --archive-mmap       (memory-mapped archive file)
    --archive-preload    (preload archive file)
    --bind <ip>:<port>
    --processes <count>  (of slave processes)
    --threads <count>    (of threads per process)
    --profile-cpu <path>
    --profile-mem <path>
    --debug




Examples
--------

* fetch and extract the Python 3.7 documentation HTML archive: ::

    curl -s -S -f \
        https://docs.python.org/3/archives/python-3.7.1-docs-html.tar.bz2 \
    | tar -x -j -v

* create the CDB archive (without any compression): ::

    cdb-http-archiver \
        --archive ./python-3.7.1-docs.cdb \
        --sources ./python-3.7.1-docs-html \
        --debug

* create the CDB archive (with `gzip` compression): ::

    cdb-http-archiver \
        --archive ./python-3.7.1-docs-gzip.cdb \
        --sources ./python-3.7.1-docs-html \
        --compress gzip \
        --debug

* serve the CDB archive (with `gzip` compression): ::

    cdb-http-server \
        --bind 127.0.0.1:8080 \
        --archive ./python-3.7.1-docs-gzip.cdb \
        --preload \
        --debug

* compare sources and archive sizes: ::

    du -h -s \
        ./python-3.7.1-docs-html \
        ./python-3.7.1-docs.cdb \
        ./python-3.7.1-docs-gzip.cdb

    46M     ./python-3.7.1-docs-html
    45M     ./python-3.7.1-docs.cdb
    9.6M    ./python-3.7.1-docs-gzip.cdb




Installation
============




Download binaries
-----------------

.. warning ::

  No binaries are currently available for download!
  Please consult the `Build from sources`_ section for now.




Build from sources
------------------


Install the prerequisites
.........................

* Ubuntu / Debian: ::

    apt-get install git-core
    apt-get install golang
    apt-get install libbrotli-dev

* OpenSUSE: ::

    zypper install git-core
    zypper install go
    zypper install libbrotli-devel


Fetch the sources
.................

::

    git clone \
        --depth 1 \
        --recurse-submodules --shallow-submodules \
        https://github.com/cipriancraciun/go-cdb-http.git \
        /tmp/go-cdb-http/src


Compile the binaries
....................

Prepare the Go environment: ::

    mkdir /tmp/go-cdb-http/go
    ln -s -T ../src/vendor /tmp/go-cdb-http/go/src

Compile the Go binnaries: ::

    export GOPATH=/tmp/go-cdb-http/go

    mkdir /tmp/go-cdb-http/bin

    go build \
        -ldflags '-s' \
        -o /tmp/go-cdb-http/bin/cdb-http-archiver \
        /tmp/go-cdb-http/src/sources/cmd/archiver.go

    go build \
        -ldflags '-s' \
        -o /tmp/go-cdb-http/bin/cdb-http-server \
        /tmp/go-cdb-http/src/sources/cmd/server.go


Deploy the binaries
...................

(Basically just copy the two executables anywhere on the system, or any compatible remote system.)

::

    cp /tmp/go-cdb-http/bin/cdb-http-archiver /usr/local/bin
    cp /tmp/go-cdb-http/bin/cdb-http-server /usr/local/bin




Features
========

The following is a list of the most important features:

* (optionally)  the static content is compressed when the CDB database is created, thus no CPU cycles are used while serving requests;

* (optionally)  the static content can be compressed with either `gzip` or Brotli_;

* (optionally)  in order to reduce the serving latency even further, one can preload the entire CDB database in memory, or alternatively mapping it in memory (mmap_);  this trades memory for CPU;

* "atomic" site content changes;  because the entire site content is held in a single CDB database file, and because the file replacement is atomically achieved via the `rename` syscall (or the `mv` tool), all the site's resources are "changed" at the same time;




Pending
=======

The following is a list of the most important features that are currently missing and are planed to be implemented:

* support for HTTPS;  (although for HTTPS it is strongly recommended to use a dedicated TLS terminator like HAProxy_;)

* support for mapping virtual hosts to multiple CDB database files;  (i.e. the ability to serve multiple domains, each with its own CDB database;)

* automatic reloading of CDB database files;

* customized error pages (also part of the CDB database);




Limitations
===========

As stated in the `About`_ section, nothing comes for free, and in order to provide all these features, some corners had to be cut:

* the CDB database **maximum size is 2 GiB**;  (however if you have a site this large, you are probabbly doing something extreemly wrong;)

* the server **does not support per-request decompression / recompression**;  this implies that if the site content was saved in the CDB database with compression (say `gzip`), the server will serve all resources compressed (i.e. `Content-Encoding : gzip`), regardless of what the browser accepts (i.e. `Accept-Encoding: gzip`);  the same applies for uncompressed content;  (however always using `gzip` compression is safe enough as it is implemented in virtually all browsers and HTTP clients out there;)

* (TODO)  currently if the CDB database file changes, the server needs to be restarted in order to pickup the changed files;

* regarding the "atomic" site changes, there is a small time window in which a client that has fetched an "old" version of a resource (say an HTML page), but which has not yet fetched the required resources (say the CSS or JS files), and the CDB database was swapped, it will consequently fetch the "new" version of these required resources;  however due to the low latency serving, this time window is extreemly small;  (**this is not a limitation of this HTTP server, but a limitation of the way the "web" is built;**)




Authors
=======

Ciprian Dorin Craciun
  * `ciprian@volution.ro <mailto:ciprian@volution.ro>`_ or `ciprian.craciun@gmail.com <mailto:ciprian.craciun@gmail.com>`_
  * `<https://volution.ro/ciprian>`_
  * `<https://github.com/cipriancraciun>`_




Notice (copyright and licensing)
================================


Notice -- short version
-----------------------

The code is licensed under AGPL 3 or later.

If you **change** the code within this repository **and use** it for **non-personal** purposes, you'll have to release it as per AGPL.


Notice -- long version
----------------------

For details about the copyright and licensing, please consult the `notice <./documentation/licensing/notice.txt>`__ file in the `documentation/licensing <./documentation/licensing>`_ folder.

If someone requires the sources and/or documentation to be released
under a different license, please send an email to the authors,
stating the licensing requirements, accompanied with the reasons
and other details; then, depending on the situation, the authors might
release the sources and/or documentation under a different license.




References
==========


.. [CDB] `CDB @WikiPedia <https://goo.gl/nvWKcY>`_

.. [Brotli] `Brotli @WikiPedia <https://goo.gl/qJHmdm>`_

.. [mmap] `Memory mapping @WikiPedia <https://goo.gl/3u6pXC>`_

.. [HAProxy] `HAProxy Load Balancer <https://goo.gl/43dnu8>`_

.. [wrk] `wrk -- modern HTTP benchmarking tool <https://goo.gl/BjpjND>`_
