

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




About
=====

``kawipiko`` is a simple static website HTTP server written in Go_, whose main purpose is to serve static website content as fast as possible.
However "simple" doesn't imply "dumb" or "limited", instead it implies "efficiency" and removal of superfluous features, inline with UNIX's philosophy of `do one thing and do it well <https://en.wikipedia.org/wiki/Unix_philosophy#Do_One_Thing_and_Do_It_Well>`__.
As such ``kawipiko`` basically supports only ``GET`` (and ``HEAD``) requests and does not provide features like dynamic content, authentication, reverse proxying, etc.

However, ``kawipiko`` does provide something unique, that no other HTTP server offers:  the static website content is served from a CDB_ database with almost zero latency.
Moreover, the static website content can be compressed (with either ``gzip`` or ``brotli``) ahead of time, thus reducing not only CPU but also bandwidth and latency.

CDB_ databases are binary files that provide efficient read-only key-value lookup tables, initially used in some DNS and SMTP servers, mainly for their low overhead lookup operations, zero locking in multi-threaded / multi-process scenarios, and "atomic" multi-record updates.
This also makes them suitable for low-latency static website content serving over HTTP, which this project provides.

For those familiar with Netlify_, ``kawipiko`` is a "host-it-yourself" alternative featuring:

* simple deployment and configuration;  (i.e. just `fetch the binaries <#installation>`__ and use the `proper flags <#kawipiko-server>`__;)
* low and constant resource consumption (both in terms of CPU and RAM);  (i.e. you won't have surprises when under load;)
* (hopefully) extremely secure;  (i.e. it doesn't launch processes, it doesn't open any files, etc.;  basically you can easily ``chroot`` it;)

For a complete list of features please consult the `features section <#features>`__.
Unfortunately, there are also some tradeoffs as described in the `limitations section <#limitations>`__ (although none are critical).

With regard to performance, as described in the `benchmarks section <#benchmarks>`__, ``kawipiko`` is on par with NGinx, sustaining 72k requests / second with 0.4ms latency for 99% of the requests even on my 6 years old laptop.
However the main advantage over NGinx is not raw performance, but deployment and configuration simplicity, plus efficient management and storage of large collections of many small files.




.. contents::
    :depth: 1
    :backlinks: none




::

    +---------------------------------------------------------------------------+
    .                                                                           .
    .   __                                                 __                   .
    .  /\ \                              __            __ /\ \                  .
    .  \ \ \/'\       __     __  __  __ /\_\   _____  /\_\\ \ \/'\      ___     .
    .   \ \ , <     /'__`\  /\ \/\ \/\ \\/\ \ /\ '__`\\/\ \\ \ , <     / __`\   .
    .    \ \ \\`\  /\ \L\.\_\ \ \_/ \_/ \\ \ \\ \ \L\ \\ \ \\ \ \\`\  /\ \L\ \  .
    .     \ \_\ \_\\ \__/.\_\\ \___x___/' \ \_\\ \ ,__/ \ \_\\ \_\ \_\\ \____/  .
    .      \/_/\/_/ \/__/\/_/ \/__//__/    \/_/ \ \ \/   \/_/ \/_/\/_/ \/___/   .
    .                                            \ \_\                          .
    .                                             \/_/                          .
    .                                                                           .
    .            _  _ ___ ___ ___     ____ ____ ____ _  _ ____ ____             .
    .            |__|  |   |  |__]    [__  |___ |__/ |  | |___ |__/             .
    .            |  |  |   |  |       ___] |___ |  \  \/  |___ |  \             .
    .                                                                           .
    .                                                                           .
    +---------------------------------------------------------------------------+




Documentation
=============

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Workflow
--------

The project provides two binaries:

* ``kawipiko-server`` -- which serves the static website content from the CDB file;
* ``kawipiko-archiver`` -- which creates the CDB file from a source folder holding the static website content;

Unlike most (if not all) other servers out-there, in which you just point your web server to the folder holding the static website content root, ``kawipiko`` takes a radically different approach.
In order to serve the static website content, one has to first "compile" it into the CDB file through ``kawipiko-archiver``, and then one can "serve" it from the CDB file through ``kawipiko-server``.

This two step phase also presents a few opportunities:

* one can decouple the "building", "testing", and "publishing" phases of a static website, by using a similar CI/CD pipeline as done for other software projects;
* one can instantaneously rollback to a previous version if the newly published one has issues;


.. note ::

   As described in the `limitations section <#limitations>`__, at the moment, if one rebuilds the CDB file, the server has to be restarted.




``kawipiko-archiver``
---------------------


::

    >> kawipiko-archiver --help

::

    Usage of kawipiko-archiver:

    --sources <path>

    --archive <path>
    --compress <gzip | brotli | identity>

    --exclude-index
    --exclude-strip
    --exclude-cache
    --include-etag

    --exclude-file-listing
    --include-folder-listing

    --debug


Flags
.....

``--sources``
    The path to the input folder that is the root of the static website content.

``--archive``
    The path to the output CDB file that contains the archived static website content.

``--compress``
    Each individual file (and consequently of the corresponding HTTP response body) is compressed with either ``gzip`` or Brotli_;  by default (or alternatively ``identity``) no compression is used.
    Even if compression is explicitly requested, if the compression ratio is bellow a certain threshold (depending on the uncompressed size), the file is stored without any compression.
    (It's senseless to force the client to spend time and decompress the response body if that time is not recovered during network transmission.)

``--exclude-index``
    Disables using ``index.*`` files (where ``.*`` is one of ``.html``, ``.htm``, ``.xhtml``, ``.xht``, ``.txt``, ``.json``, and ``.xml``) to respond to a request whose URL ends in ``/`` (corresponding to the folder wherein ``index.*`` file is located).
    (This can be used to implement "slash" blog style URL's like ``/blog/whatever/`` which maps to ``/blog/whatever/index.html``.)

``--exclude-strip``
    Disables using a file with the suffix ``.html``, ``.htm``, ``.xhtml``, ``.xht``, and ``.txt`` to respond to a request whose URL does not exactly match an existing file.
    (This can be used to implement "suffix-less" blog style URL's like ``/blog/whatever`` which maps to ``/blog/whatever.html``.)

``--exclude-cache``
    Disables adding an ``Cache-Control: public, immutable, max-age=3600`` header that forces the browser (and other intermediary proxies) to cache the response for an hour (the ``public`` and ``max-age=3600`` arguments), and furthermore not request it even on reloads (the ``immutable`` argument).

``--include-etag``
    Enables adding an ``ETag`` response header that contains the SHA256 of the response body.
    By not including the ``ETag`` header (i.e. the default), and because identical headers are stored only one, if one has many files of the same type (that in turn without ``ETag`` generates the same headers), this can lead to significant reduction in stored headers, including reducing RAM usage.
    (At this moment it does not support HTTP conditional requests, i.e. the ``If-None-Match``, ``If-Modified-Since`` and their counterparts;  however this ``ETag`` header might be used in conjuction with ``HEAD`` requests to see if the resource has changed.)

``--exclude-file-listing``
    Disables the creation of an internal list of files that can be used in conjunction with the ``--index-all`` flag of the ``kawipiko-server``.

``--include-folder-listing``
    Enables the creation of an internal list of folders.  (Currently not used by the ``kawipiko-server`` tool.)

``--debug``
    Enables verbose logging.
    It will log various information about the archived files (including compression statistics).


Ignored files
.............

* any file with the following prefixes: ``.``, ``#``;
* any file with the following suffixes: ``~``, ``#``, ``.log``, ``.tmp``, ``.temp``, ``.lock``;
* any file that contains the following: ``#``;
* any file that exactly matches the following:: ``Thumbs.db``, ``.DS_Store``;
* (at the moment these rules are not configurable through flags;)


``_wildcard.*`` files
.....................


By placing a file whose name matches ``_wildcard.*`` (i.e. with the prefix ``_wildcard.`` and any other suffix), it will be used to respond to any request whose URL fails to find a "better" match.

These wildcard files respect the folder hierarchy, in that wildcard files in (direct or transitive) subfolders override the wildcard file in their parents (direct or transitive).


Symlinks, hardlinks, loops, and duplicated files
................................................

You freely use symlinks (including pointing outside of the content root) and they will be crawled during archival respecting the "logical" hierarchy they introduce.
(Any loop that you introduce into the hierarchy will be ignored and a warning will be issued.)

You can safely symlink or hardlink the same file (or folder) in multiple places (within the content hierarchy), and its data will be stored only once.
(The same applies to duplicated files that have exactly the same data.)




``kawipiko-server``
-------------------

::

    >> kawipiko-server --help

::

    Usage of kawipiko-server:

    --archive <path>
    --archive-inmem      (memory-loaded archive file)
    --archive-mmap       (memory-mapped archive file)
    --archive-preload    (preload archive file)

    --index-all
    --index-paths
    --index-data-meta
    --index-data-content

    --bind <ip>:<port>

    --processes <count>  (of slave processes)
    --threads <count>    (of threads per process)

    --security-headers-tls
    --security-headers-disable

    --profile-cpu <path>
    --profile-mem <path>

    --debug
    --dummy


Flags
.....


``--archive``
    The path of the CDB file that contains the archived static website content.
    (It can be created with the ``kawipiko-archiver`` tool.)

``--archive-inmem``
    Reads the CDB file in memory, and thus all requests are served from RAM.
    (This can be used if enough RAM is available to avoid swapping.)

``--archive-mmap``
    The CDB file is `memory mapped <#mmap>`__.
    (**Highly recommended!**)

``--archive-preload``
    Before starting to serve requests, read the CDB file so that its data is buffered by the OS.
    (**Highly recommended!**)

``--index-all``, ``--index-paths``, ``--index-data-meta``,  and ``--index-data-content``
    In order to serve a request:

    * the request URL's path is used to locate a resource's metadata (i.e. response headers) and data (i.e. response body) fingerprints;
      by using ``--index-paths`` an RAM-based hash-map is created to eliminate a CDB lookup operation for this purpose;

    * based on the resource's metadata fingerprint, the actual metadata (i.e. the response headers) is located;
      by using ``--index-data-meta`` a RAM-based hash-map is created to eliminate a CDB lookup operation for this purpose;

    * based on the resource's data fingerprint, the actual data (i.e. the response body) is located;
      by using ``--index-data-content`` a RAM-based hash-map is created to eliminate a CDB lookup operation for this purpose;

    * ``--index-all`` enables all these indices;

    * (depending on the use-case) it is highly recommended to use ``--index-paths``;   if ``--exclude-etag`` was used during archival, one can also use ``--index-data-meta``;

    * it is highly recommended to use ``--archive-inmem`` or ``--archive-mmap`` or else (especially if data is indexed) the net effect is that of loading everything in RAM;

``--bind``
    The IP and port to listen for requests.

``--processes`` and ``--threads``
    The number of processes and threads per each process to start.
    It is highly recommended to use 1 process and as many threads as there are cores.

    Depending on the use-case, one can use multiple processes each with a single thread;  this would reduce goroutine contention if it causes problems.
    (However note that if using ``--archive-inmem`` each process will allocate its own copy of the database in RAM;  in such cases it is highly recommended to use ``--archive-mmap``.)

``--security-headers-tls``
    Enables adding the ``Strict-Transport-Security: max-age=31536000`` and ``Content-Security-Policy: upgrade-insecure-requests`` to the response headers.
    (Although at the moment ``kawipiko`` does not support HTTPS, it can be used behind a TLS terminator, load-balancer or proxy that do support HTTPS;  therefore these headers instruct the browser to always use HTTPS for the served domain.)

``--security-headers-disable``
    Disables adding a few security related headers: ::

      Referrer-Policy: strict-origin-when-cross-origin
      X-Content-Type-Options: nosniff
      X-XSS-Protection: 1; mode=block
      X-Frame-Options: sameorigin

``--debug``
    Enables verbose logging.
    (**Highly discouraged!**)

``--dummy``
    It starts the server in "dummy" mode, ignoring all archive related arguments and always responding with ``hello world!\n`` and without additional headers except the HTTP status line and ``Content-Length``.
    This argument can be used to benchmark the raw performance of the underlying Go and ``fasthttp`` performance;  this is the upper limit on the achievable performance given the underlying technologies.
    (From my own benchmarks ``kawipiko``'s adds only about ~15% overhead when actually serving the ``hello-world.cdb`` archive.)

``--profile-cpu`` and `--profile-mem``
    Enables CPU and memory profiling using Go's profiling infrastructure.




Examples
--------

* fetch and extract the Python 3.7 documentation HTML archive: ::

    curl -s -S -f \
        -o ./python-3.7.3-docs-html.tar.bz2 \
        https://docs.python.org/3/archives/python-3.7.3-docs-html.tar.bz2 \
    #

    tar -x -j -v -f ./python-3.7.3-docs-html.tar.bz2

* create the CDB archive (without any compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.3-docs-html-nozip.cdb \
        --sources ./python-3.7.3-docs-html \
        --debug \
    #

* create the CDB archive (with ``gzip`` compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.3-docs-html-gzip.cdb \
        --sources ./python-3.7.3-docs-html \
        --compress gzip \
        --debug \
    #

* create the CDB archive (with ``brotli`` compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.3-docs-html-brotli.cdb \
        --sources ./python-3.7.3-docs-html \
        --compress brotli \
        --debug \
    #

* serve the CDB archive (with ``gzip`` compression): ::

    kawipiko-server \
        --bind 127.0.0.1:8080 \
        --archive ./python-3.7.3-docs-html-gzip.cdb \
        --archive-mmap \
        --archive-preload \
        --debug \
    #

* compare sources and archive sizes: ::

    du -h -s \
        \
        ./python-3.7.3-docs-html-nozip.cdb \
        ./python-3.7.3-docs-html-gzip.cdb \
        ./python-3.7.3-docs-html-brotli.cdb \
        \
        ./python-3.7.3-docs-html \
        ./python-3.7.3-docs-html.tar.bz2 \
    #

    45M     ./python-3.7.3-docs-html-nozip.cdb
    9.7M    ./python-3.7.3-docs-html-gzip.cdb
    7.9M    ./python-3.7.3-docs-html-brotli.cdb

    46M     ./python-3.7.3-docs-html
    6.0M    ./python-3.7.3-docs-html.tar.bz2




Installation
============

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Download binaries
-----------------

.. warning ::

  No binaries are currently available for download!
  Please consult the `build from sources section <#build-from-sources>`__ for now.




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
        https://github.com/volution/kawipiko.git \
        /tmp/kawipiko/src \
    #


Compile the binaries
....................

Prepare the Go environment: ::

    export GOPATH=/tmp/kawipiko/go

    mkdir /tmp/kawipiko/go
    mkdir /tmp/kawipiko/bin

Compile the Go binaries: ::

    cd /tmp/kawipiko/src/sources

    go build \
        -ldflags '-s' \
        -o /tmp/kawipiko/bin/kawipiko-archiver \
        ./cmd/archiver.go \
    #

    go build \
        -ldflags '-s' \
        -o /tmp/kawipiko/bin/kawipiko-server \
        ./cmd/server.go \
    #


Deploy the binaries
...................

(Basically just copy the two executables anywhere on the system, or any compatible remote system.)

::

    cp /tmp/kawipiko/bin/kawipiko-archiver /usr/local/bin
    cp /tmp/kawipiko/bin/kawipiko-server /usr/local/bin




Features
========

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Implemented
-----------

The following is a list of the most important features:

* (optionally)  the static website content is compressed when the CDB database is created, thus no CPU cycles are used while serving requests;

* (optionally)  the static website content can be compressed with either ``gzip`` or Brotli_;

* (optionally)  in order to reduce the serving latency even further, one can preload the entire CDB database in memory, or alternatively mapping it in memory (mmap_);  this trades memory for CPU;

* "atomic" static website content changes;  because the entire content is held in a single CDB database file, and because the file replacement is atomically achieved via the ``rename`` syscall (or the ``mv`` tool), all resources are "changed" at the same time;

* ``_wildcard.*`` files (where ``.*`` are the regular extensions like ``.txt``, ``.html``, etc.) which will be used if an actual resource is not found under that folder;  (these files respect the hierarchical tree structure, i.e. "deeper" ones override the ones closer to "root";)




Pending
-------

The following is a list of the most important features that are currently missing and are planed to be implemented:

* support for HTTPS;  (although for HTTPS it is strongly recommended to use a dedicated TLS terminator like HAProxy_;)

* support for custom HTTP response headers (for specific files, for specific folders, etc.);  (currently only ``Content-Type``, ``Content-Length``, ``Content-Encoding`` and optionally ``ETag`` is included;  additionally ``Cache-Control: public, immutable, max-age=3600`` and a few security related headers are also included;)

* support for mapping virtual hosts to key prefixes;  (currently virtual hosts, i.e. the ``Host`` header, are ignored;)

* support for mapping virtual hosts to multiple CDB database files;  (i.e. the ability to serve multiple domains, each with its own CDB database;)

* automatic reloading of CDB database files;

* customized error pages (also part of the CDB database);




Limitations
-----------

As stated in the `about section <#about>`__, nothing comes for free, and in order to provide all these features, some corners had to be cut:

* (TODO)  currently if the CDB database file changes, the server needs to be restarted in order to pickup the changed files;

* (won't fix)  the CDB database **maximum size is 4 GiB**;  (however if you have a static website this large, you are probably doing something extremely wrong, as large files should be offloaded to something like AWS S3 and served through a CDN like CloudFlare or AWS CloudFront;)

* (won't fix)  the server **does not support per-request decompression / recompression**;  this implies that if the content was saved in the CDB database with compression (say ``gzip``), the server will serve all resources compressed (i.e. ``Content-Encoding: gzip``), regardless of what the browser accepts (i.e. ``Accept-Encoding: gzip``);  the same applies for uncompressed content;  (however always using ``gzip`` compression is safe enough as it is implemented in virtually all browsers and HTTP clients out there;)

* (won't fix)  regarding the "atomic" static website changes, there is a small time window in which a client that has fetched an "old" version of a resource (say an HTML page), but which has not yet fetched the required resources (say the CSS or JS files), and the CDB database was swapped, it will consequently fetch the "new" version of these required resources;  however due to the low latency serving, this time window is extremely small;  (**this is not a limitation of this HTTP server, but a limitation of the way the "web" is built;**  always use fingerprints in your resources URL, and perhaps always include the current and previous version on each deploy;)




Benchmarks
==========

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Summary
-------

Bottom line (**even on my 6 years old laptop**):

* under normal conditions (16 concurrent connections), you get around 111k requests / second, at about 0.25ms latency for 99% of the requests;
* under light stress conditions (128 concurrent connections), you get around 118k requests / second, at about 2.5ms latency for 99% of the requests;
* under medium stress conditions (512 concurrent connections), you get around 106k requests / second, at about 10ms latency for 99% of the requests (meanwhile the average is 4.5ms);
* **under high stress conditions (2048 concurrent connections), you get around 100k requests / second, at about 400ms latency for 99% of the requests (meanwhile the average is 45ms);**
* under extreme stress conditions (16384 concurrent connections) (i.e. someone tries to DDOS the server), you get around 53k requests / second, at about 2.8s latency for 99% of the requests (meanwhile the average is 200ms);
* (the timeout errors are due to the fact that ``wrk`` is configured to timeout after only 1 second of waiting while connecting or receiving the full response;)
* (the read errors are due to the fact that the server closes a keep-alive connection after serving 256k requests;)
* **the raw performance is comparable with NGinx_** (only 20% few requests / second for this "synthetic" benchmark);  however for a "real" scenario (i.e. thousand of small files accessed in a random pattern) I think they are on-par;  (not to mention how simple it is to configure and deploy ``kawipiko`` as compared to NGinx;)




Results
-------


Results values
..............


.. note ::

  Please note that the values under *Thread Stats* are reported per thread.
  Therefore it is best to look at the first two values, i.e. *Requests/sec*.

* 16 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec: 111720.73
    Transfer/sec:     18.01MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 16 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   139.36us   60.27us   1.88ms   64.91%
        Req/Sec    56.14k   713.04    57.60k    91.36%
      Latency Distribution
         50%  143.00us
         75%  184.00us
         90%  212.00us
         99%  261.00us
      3362742 requests in 30.10s, 541.98MB read

* 128 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec: 118811.41
    Transfer/sec:     19.15MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 128 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     1.03ms  705.69us  19.53ms   63.54%
        Req/Sec    59.71k     1.69k   61.70k    96.67%
      Latency Distribution
         50%    0.99ms
         75%    1.58ms
         90%    1.89ms
         99%    2.42ms
      3564527 requests in 30.00s, 574.50MB read

* 512 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec: 106698.89
    Transfer/sec:     17.20MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     4.73ms    3.89ms  39.32ms   39.74%
        Req/Sec    53.71k     1.73k   69.18k    84.33%
      Latency Distribution
         50%    4.96ms
         75%    8.63ms
         90%    9.19ms
         99%   10.30ms
      3206540 requests in 30.05s, 516.80MB read
      Socket errors: connect 0, read 105, write 0, timeout 0

* 2048 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec: 100296.65
    Transfer/sec:     16.16MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    45.42ms   85.14ms 987.70ms   88.62%
        Req/Sec    50.61k     5.59k   70.14k    71.74%
      Latency Distribution
         50%   16.30ms
         75%   28.44ms
         90%  147.60ms
         99%  417.40ms
      3015868 requests in 30.07s, 486.07MB read
      Socket errors: connect 0, read 128, write 0, timeout 86

* 4096 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec:  95628.34
    Transfer/sec:     15.41MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 4096 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    90.50ms  146.08ms 999.65ms   88.49%
        Req/Sec    48.27k     6.09k   66.05k    76.34%
      Latency Distribution
         50%   23.31ms
         75%  112.06ms
         90%  249.41ms
         99%  745.94ms
      2871404 requests in 30.03s, 462.79MB read
      Socket errors: connect 0, read 27, write 0, timeout 4449

* 16384 connections / 2 server threads / 2 wrk threads: ::

    Requests/sec:  53548.52
    Transfer/sec:      8.63MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   206.21ms  513.75ms   6.00s    92.56%
        Req/Sec    31.37k     5.68k   44.44k    76.13%
      Latency Distribution
         50%   35.38ms
         75%   62.78ms
         90%  551.33ms
         99%    2.82s
      1611294 requests in 30.09s, 259.69MB read
      Socket errors: connect 0, read 115, write 0, timeout 2288


Results notes
.............

* the machine was my personal laptop:  6 years old with an Intel Core i7 3667U (2 cores with 2 threads each);
* the ``kawipiko-server`` was started with ``--processes 1 --threads 2``;  (i.e. 2 threads handling the requests;)
* the ``kawipiko-server`` was started with ``--archive-inmem``;  (i.e. the CDB database file was preloaded into memory, thus no disk I/O;)
* the ``kawipiko-server`` was started with ``--security-headers-disable``;  (because these headers are not set by default by other HTTP servers;)
* the ``kawipiko-server`` was started with ``--timeout-disable``;  (because, due to a known Go issue, using ``net.Conn.SetDeadline`` has an impact of about 20% of the raw performance;  thus the reported values above might be about 10%-15% smaller when used with timeouts;)
* the benchmarking tool was wrk_;
* both ``kawipiko-server`` and ``wrk`` tools were run on the same machine;
* both ``kawipiko-server`` and ``wrk`` tools were pinned on different physical cores;
* the benchmark was run over loopback networking (i.e. ``127.0.0.1``);
* the served file contains ``Hello World!``;
* the protocol was HTTP (i.e. no TLS), with keep-alive;
* see the `methodology section <#methodology>`__ for details;




Comparisons
-----------


Comparisons with NGinx
......................

* NGinx 512 connections / 2 server workers / 2 wrk thread: ::

    Requests/sec:  79816.08
    Transfer/sec:     20.02MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     6.07ms    1.90ms  19.83ms   71.67%
        Req/Sec    40.17k     1.16k   43.35k    69.83%
      Latency Distribution
         50%    6.13ms
         75%    6.99ms
         90%    8.51ms
         99%   11.10ms
      2399069 requests in 30.06s, 601.73MB read

* NGinx 2048 connections / 2 server workers / 2 wrk thread: ::

    Requests/sec:  78211.46
    Transfer/sec:     19.62MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    27.11ms   20.27ms 490.12ms   97.76%
        Req/Sec    39.45k     2.45k   49.98k    70.74%
      Latency Distribution
         50%   24.80ms
         75%   29.67ms
         90%   34.99ms
         99%  126.97ms
      2351933 requests in 30.07s, 589.90MB read
      Socket errors: connect 0, read 0, write 0, timeout 11

* NGinx 4096 connections / 2 server workers / 2 wrk thread: ::

    Requests/sec:  75970.82
    Transfer/sec:     19.05MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 4096 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    70.25ms   73.68ms 943.82ms   87.21%
        Req/Sec    38.37k     3.79k   49.06k    70.30%
      Latency Distribution
         50%   46.37ms
         75%   58.28ms
         90%  179.08ms
         99%  339.05ms
      2282223 requests in 30.04s, 572.42MB read
      Socket errors: connect 0, read 0, write 0, timeout 187

* NGinx 16384 connections / 2 server workers / 2 wrk thread: ::

    Requests/sec:  43909.67
    Transfer/sec:     11.01MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   223.87ms  551.14ms   5.94s    92.92%
        Req/Sec    32.95k    13.35k   51.56k    76.71%
      Latency Distribution
         50%   32.62ms
         75%  222.93ms
         90%  558.04ms
         99%    3.17s
      1320562 requests in 30.07s, 331.22MB read
      Socket errors: connect 0, read 12596, write 34, timeout 1121

* (the NGinx configuration file can be found in the `examples folder <./examples>`__;  the configuration was obtained after many experiments to squeeze out of NGinx as much performance as possible, given the targeted use-case, namely many small files;)


Comparisons with others
.......................

* darkhttpd_ 512 connections / 1 server process / 2 wrk threads: ::

    Requests/sec:  38191.65
    Transfer/sec:      8.74MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    17.51ms   17.30ms 223.22ms   78.55%
        Req/Sec     9.62k     1.94k   17.01k    72.98%
      Latency Distribution
         50%    7.51ms
         75%   32.51ms
         90%   45.69ms
         99%   53.00ms
      1148067 requests in 30.06s, 262.85MB read




Methodology
-----------


* get the binaries (either `download <#download-binaries>`__ or `build <#build-from-sources>`__ them);
* get the ``hello-world.cdb`` (from the `examples <./examples>`__ folder inside the repository);


Single process / single threaded
................................

* this scenario will yield a "base-line performance" per core;

* execute the server (in-memory and indexed) (i.e. the "best case scenario"): ::

    kawipiko-server \
        --bind 127.0.0.1:8080 \
        --archive ./hello-world.cdb \
        --archive-inmem \
        --index-all \
        --processes 1 \
        --threads 1 \
    #

* execute the server (memory mapped) (i.e. the "the recommended scenario"): ::

    kawipiko-server \
        --bind 127.0.0.1:8080 \
        --archive ./hello-world.cdb \
        --archive-mmap \
        --processes 1 \
        --threads 1 \
    #


Single process / two threads
............................

* this scenario is the usual setup;  configure ``--threads`` to equal the number of cores;

* execute the server (memory mapped): ::

    kawipiko-server \
        --bind 127.0.0.1:8080 \
        --archive ./hello-world.cdb \
        --archive-mmap \
        --processes 1 \
        --threads 2 \
    #


Load generators
...............

* 512 concurrent connections (handled by 2 threads): ::

    wrk \
        --threads 2 \
        --connections 512 \
        --timeout 1s \
        --duration 30s \
        --latency \
        http://127.0.0.1:8080/ \
    #

* 4096 concurrent connections (handled by 2 threads): ::

    wrk \
        --threads 2 \
        --connections 4096 \
        --timeout 1s \
        --duration 30s \
        --latency \
        http://127.0.0.1:8080/ \
    #


Methodology notes
.................

* the number of threads for the server plus for ``wkr`` shouldn't be larger than the number of available cores;  (or use different machines for the server and the client;)

* also take into account that by default the number of "file descriptors" on most UNIX/Linux machines is 1024, therefore if you want to try with more connections than 1000, you need to raise this limit;  (see bellow;)

* additionally, you can try to pin the server and ``wrk`` to specific cores, increase various priorities (scheduling, IO, etc.);  (given that Intel processors have HyperThreading which appear to the OS as individual cores, you should make sure that you pin each process on cores part of the same physical processor / core;)

* pinning the server (cores ``0`` and ``1`` are mapped on physical core ``1``): ::

    sudo -u root -n -E -P -- \
    \
    taskset -c 0,1 \
    nice -n -19 -- \
    ionice -c 2 -n 0 -- \
    chrt -r 10 \
    prlimit -n262144 -- \
    \
    sudo -u "${USER}" -n -E -P -- \
    \
    kawipiko-server \
        ... \
    #

* pinning the client (cores ``2`` and ``3`` are mapped on physical core ``2``): ::

    sudo -u root -n -E -P -- \
    \
    taskset -c 2,3 \
    nice -n -19 -- \
    ionice -c 2 -n 0 -- \
    chrt -r 10 \
    prlimit -n262144 -- \
    \
    sudo -u "${USER}" -n -E -P -- \
    \
    wrk \
        ... \
    #




Why CDB?
========

Until I expand upon why I have chosen to use CDB for service static website content, you can read about the `sparkey <https://github.com/spotify/sparkey>`__ from Spotify.




Authors
=======

Ciprian Dorin Craciun
  * `ciprian@volution.ro <mailto:ciprian@volution.ro>`__ or `ciprian.craciun@gmail.com <mailto:ciprian.craciun@gmail.com>`__
  * `<https://volution.ro/ciprian>`__
  * `<https://github.com/cipriancraciun>`__




Notice (copyright and licensing)
================================

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Notice -- short version
-----------------------

The code is licensed under AGPL 3 or later.

If you **change** the code within this repository **and use** it for **non-personal** purposes, you'll have to release it as per AGPL.




Notice -- long version
----------------------

For details about the copyright and licensing, please consult the `notice <./documentation/licensing/notice.txt>`__ file in the `documentation/licensing <./documentation/licensing>`__ folder.

If someone requires the sources and/or documentation to be released
under a different license, please send an email to the authors,
stating the licensing requirements, accompanied with the reasons
and other details; then, depending on the situation, the authors might
release the sources and/or documentation under a different license.




References
==========


.. [Go]
    * `Go <https://en.wikipedia.org/wiki/Go_(programming_language)>`__ (@WikiPedia);
    * `Go <https://golang.com/>`__ (project);

.. [CDB]
    * `CDB <https://en.wikipedia.org/wiki/Cdb_(software)>`__ (@WikiPedia);
    * `cdb <http://cr.yp.to/cdb.html>`__ (project);
    * `cdb internals <http://www.unixuser.org/~euske/doc/cdbinternals/index.html>`__ (article);
    * `Benchmarking LevelDB vs. RocksDB vs. HyperLevelDB vs. LMDB Performance for InfluxDB <https://www.influxdata.com/blog/benchmarking-leveldb-vs-rocksdb-vs-hyperleveldb-vs-lmdb-performance-for-influxdb/>`__ (article);
    * `Badger vs LMDB vs BoltDB: Benchmarking key-value databases in Go <https://blog.dgraph.io/post/badger-lmdb-boltdb/>`__ (article);
    * `Benchmarking BDB, CDB and Tokyo Cabinet on large datasets <https://www.dmo.ca/blog/benchmarking-hash-databases-on-large-data/>`__ (article);
    * `TinyCDB <http://www.corpit.ru/mjt/tinycdb.html>`__ (fork project);
    * `tinydns <https://cr.yp.to/djbdns/tinydns.html>`__ (DNS server using CDB);
    * `qmail <https://cr.yp.to/qmail.html>`__ (SMTP server using CDB);

.. [wrk]
    * `wrk <https://github.com/wg/wrk>`__ (project);
    * modern HTTP benchmarking tool;
    * multi threaded, with event loop and Lua support;

.. [Brotli]
    * `Brotli <https://en.wikipedia.org/wiki/Brotli>`__ (@WikiPedia);
    * `Brotli <https://github.com/google/brotli>`__ (project);
    * `Results of experimenting with Brotli for dynamic web content <https://blog.cloudflare.com/results-experimenting-brotli/>`__ (article);

.. [Netlify]
    * `Netlify <https://www.netlify.com/>`__ (cloud provider);

.. [HAProxy]
    * `HAProxy <https://en.wikipedia.org/wiki/HAProxy>`__ (@WikiPedia);
    * `HAProxy <https://www.haproxy.org/>`__ (project);
    * reliable high performance TCP/HTTP load-balancer;
    * multi threaded, with event loop and Lua support;

.. [NGinx]
    * `NGinx <https://en.wikipedia.org/wiki/Nginx>`__ (@WikiPedia);
    * `NGinx <https://nginx.org/>`__ (project);

.. [darkhttpd]
    * `darkhttpd <https://unix4lyfe.org/darkhttpd/>`__ (project);
    * simple static HTTP server;
    * single threaded, with event loop and ``sendfile`` support;

.. [mmap]
    * `Memory mapping <https://en.wikipedia.org/wiki/Memory-mapped_file>`__ (@WikiPedia);
    * `mmap(2) <http://man7.org/linux/man-pages/man2/mmap.2.html>`__ (Linux man page);

