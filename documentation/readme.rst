

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




About
=====

``kawipiko`` is a simple static HTTP server written in Go_, whose main purpose is to serve static content as fast as possible.
However "simple" doesn't imply "dumb" or "limited", instead it implies "efficiency" and removal of superfluous features, inline with UNIX's philosophy of `do one thing and do it well <https://en.wikipedia.org/wiki/Unix_philosophy#Do_One_Thing_and_Do_It_Well>`__.
As such ``kawipiko`` basically supports only ``GET`` (and ``HEAD``) requests and does not provide features like dynamic content, authentication, reverse proxying, etc.

However, ``kawipiko`` does provide something unique, that no other HTTP server offers:  the static content is served from a CDB_ database with almost zero latency.
Moreover, the static content can be compressed (with either ``gzip`` or ``brotli``) ahead of time, thus reducing not only CPU but also bandwith and latency.

CDB_ databases are binary files that provide efficient read-only key-value lookup tables, initially used in some DNS and SMTP servers, mainly for their low overhead lookup operations, zero locking in multi-threaded / multi-process scenarios, and "atomic" multi-record updates.
This also makes them suitable for low-latency static content serving over HTTP, which this project provides.

For those familiar with Netlify_, ``kawipiko`` is a "host-it-yourself" alternative featuring:

* simple deployment and configuration;  (i.e. just `fetch the binaries <#installation>`__ and use the `proper flags <#kawipiko-server>`__;)
* low and constant resource consumption (both in terms of CPU and RAM);  (i.e. you won't have surprises when under load;)
* (hopefully) extremly secure;  (i.e. it doesn't launch processes, it doesn't open any files, etc.;  basically you can easly ``chroot`` it;)

For a complete list of features please consult the `features section <#features>`__.
Unfortunately, there are also some tradeoffs as described in the `limitations section <#limitations>`__ (although none are critical).

With regard to performance, as described in the `benchmarks section <#benchmarks>`__, ``kawipiko`` is on par with NGinx, sustaining 72k requests / second with 0.4ms latency for 99% of the requests even on my 6 years old laptop.
However the main advantage over NGinx is not raw performance, but deployment and configuration simplicity, plus efficient management and storage of large collections of many small static files.




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


The project provides two binaries:

* ``kawipiko-server`` -- which serves the static content from the CDB database;
* ``kawipiko-archiver`` -- which creates the CDB database from a source folder holding the static content;




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
    --exclude-etag

    --exclude-file-listing
    --include-folder-listing

    --debug


Flags
.....

``--sources``
    The path to the input folder that is the root of the website.

``--archive``
    The path to the output CDB file that contains the archived website.

``--compress``
    Each individual file (and consequently of the corresponding HTTP response body) is compressed with either ``gzip`` or Brotli_;  by default (or alternatively ``identity``) no compression is used.
    Even if compression is explicitly requested, if the compression ratio is bellow a certain threshold (depending on the uncompressed size), the file is stored without any compression.
    (It's senseless to force the client to spend time and decompress the response body if that time is not recovered during network transmission.)

``--exclude-index``
    Disables using ``index.*`` files (where ``.*`` is one of ``.html``, ``.htm``, ``.xhtml``, ``.xht``, ``.txt``, ``.json``, and ``.xml``) to respond to a request whose URL ends in ``/`` (corresponding to the folder wherein ``index.*`` file is located).
    (This can be used to implement "slash" blog style URL's like ``/blog/whatever/`` which maps to ``/blog/whatever/index.html``.)

``--exclude-strip``
    Disables using a file with the suffix ``.html``, ``.htm``, ``.xhtml``, ``.xht``, and ``.txt`` to respond to a rquest whose URL does not exactly match an existing file.
    (This can be used to implement "suffix-less" blog style URL's like ``/blog/whatever`` which maps to ``/blog/whatever.html``.)

``--exclude-etag``
    Disables adding an ``ETag`` response header that contains the SHA256 of the response body.
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

You freely use symlinks (including pointing outside of the website root) and they will be crawled during archival respecting the "logical" hierarchy they introduce.
(Any loop that you introduce into the hierarchy will be ignored and a warning will be issued.)

You can safely symlink or hardlink the same file (or folder) in multiple places (within the website root), and its contents will be stored only once.
(The same applies to duplicated files that have exactly the same contents.)




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

    --profile-cpu <path>
    --profile-mem <path>

    --debug


Flags
.....


``--archive``
    The path of the CDB file that contains the archived website.
    (It can be created with the ``kawipiko-archiver`` tool.)

``--archive-inmem``
    Reads the CDB file in memory, and thus all requests are served from RAM.
    (This can be used if enough RAM is available to avoid swapping.)

``--archive-mmap``
    The CDB file is `memory mapped <#mmap>`__.
    (**Highly recommended!**)

``--archive-preload``
    Before starting to serve requests, read the CDB file so that its contents is buffered by the OS.
    (**Highly recommended!**)

``--index-all``, ``--index-paths``, ``--index-data-meta``,  and ``--index-data-content``
    In order to serve a request:

    * the request URL's path is used to locate a resource's metadata (i.e. response headers) and data (i.e. response body) fingerprints;
      by using ``--index-paths`` an RAM-based hash-map is created to eliminate a CDB lookup operation for this purpose;

    * based on the resource's metadata fingerprint, the actual metadata (i.e. the response headers) is located;
      by using ``--index-data-meta`` a RAM-based hash-map is created to eliminate a CDB lookup operation for this purpose;

    * based on the resource's data fingerprint, the actual contents (i.e. the response body) is located;
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

``--debug``
    Enables verbose logging.
    (**Highly discouraged!**)

``--profile-cpu`` and `--profile-mem``
    Enables CPU and memory profiling using Go's profiling infrastructure.




Examples
--------

* fetch and extract the Python 3.7 documentation HTML archive: ::

    curl -s -S -f \
        -o ./python-3.7.1-docs-html.tar.bz2 \
        https://docs.python.org/3/archives/python-3.7.1-docs-html.tar.bz2 \
    #

    tar -x -j -v -f ./python-3.7.1-docs-html.tar.bz2

* create the CDB archive (without any compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.1-docs-html.cdb \
        --sources ./python-3.7.1-docs-html \
        --debug \
    #

* create the CDB archive (with ``gzip`` compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.1-docs-html-gzip.cdb \
        --sources ./python-3.7.1-docs-html \
        --compress gzip \
        --debug \
    #

* create the CDB archive (with ``brotli`` compression): ::

    kawipiko-archiver \
        --archive ./python-3.7.1-docs-html-brotli.cdb \
        --sources ./python-3.7.1-docs-html \
        --compress brotli \
        --debug \
    #

* serve the CDB archive (with ``gzip`` compression): ::

    kawipiko-server \
        --bind 127.0.0.1:8080 \
        --archive ./python-3.7.1-docs-html-gzip.cdb \
        --archive-mmap \
        --archive-preload \
        --debug \
    #

* compare sources and archive sizes: ::

    du -h -s \
        \
        ./python-3.7.1-docs-html.cdb \
        ./python-3.7.1-docs-html-gzip.cdb \
        ./python-3.7.1-docs-html-brotli.cdb \
        \
        ./python-3.7.1-docs-html \
        ./python-3.7.1-docs-html.tar.bz2 \
    #

    45M     ./python-3.7.1-docs-html.cdb
    9.9M    ./python-3.7.1-docs-html-gzip.cdb
    8.0M    ./python-3.7.1-docs-html-brotli.cdb

    46M     ./python-3.7.1-docs-html
    6.0M    ./python-3.7.1-docs-html.tar.bz2




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

Compile the Go binnaries: ::

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

* (optionally)  the static content is compressed when the CDB database is created, thus no CPU cycles are used while serving requests;

* (optionally)  the static content can be compressed with either ``gzip`` or Brotli_;

* (optionally)  in order to reduce the serving latency even further, one can preload the entire CDB database in memory, or alternatively mapping it in memory (mmap_);  this trades memory for CPU;

* "atomic" site content changes;  because the entire site content is held in a single CDB database file, and because the file replacement is atomically achieved via the ``rename`` syscall (or the ``mv`` tool), all the site's resources are "changed" at the same time;

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

* (won't fix)  the CDB database **maximum size is 4 GiB**;  (however if you have a site this large, you are probabbly doing something extreemly wrong, as large files should be offloaded to something like AWS S3 and served through a CDN like CloudFlare or AWS CloudFront;)

* (won't fix)  the server **does not support per-request decompression / recompression**;  this implies that if the site content was saved in the CDB database with compression (say ``gzip``), the server will serve all resources compressed (i.e. ``Content-Encoding: gzip``), regardless of what the browser accepts (i.e. ``Accept-Encoding: gzip``);  the same applies for uncompressed content;  (however always using ``gzip`` compression is safe enough as it is implemented in virtually all browsers and HTTP clients out there;)

* (won't fix)  regarding the "atomic" site changes, there is a small time window in which a client that has fetched an "old" version of a resource (say an HTML page), but which has not yet fetched the required resources (say the CSS or JS files), and the CDB database was swapped, it will consequently fetch the "new" version of these required resources;  however due to the low latency serving, this time window is extreemly small;  (**this is not a limitation of this HTTP server, but a limitation of the way the "web" is built;**  always use fingerprints in your resources URL, and perhaps always include the current and previous version on each deploy;)




Benchmarks
==========

.. contents::
    :depth: 2
    :local:
    :backlinks: none




Summary
-------

Bottom line (**even on my 6 years old laptop**):

* under normal conditions (16 concurrent connections), you get around 72k requests / second, at about 0.4ms latency for 99% of the requests;
* under stress conditions (512 concurrent connections), you get arround 74k requests / second, at about 15ms latency for 99% of the requests;
* **under extreme conditions (2048 concurrent connections), you get arround 74k requests / second, at about 500ms latency for 99% of the requests (meanwhile the average is 50ms);**
* (the timeout errors are due to the fact that ``wrk`` is configured to timeout after only 1 second of waiting;)
* (the read errors are due to the fact that the server closes a keep-alive connection after serving 256k requests;)
* **the raw performance is comparable with NGinx_** (only 20% few requests / second for this "synthetic" benchmark);  however for a "real" scenario (i.e. thousand of small files accessed in a random pattern) I think they are on-par;  (not to mention how simple it is to configure and deploy ``kawipiko`` as compared to NGinx;)




Results
-------


Results values
..............


.. note ::

  Please note that the values under *Thread Stats* are reported per thread.
  Therefore it is best to look at the first two values, i.e. *Requests/sec*.

* 16 connections / 2 server threads / 4 wrk threads: ::

    Requests/sec:  71935.39
    Transfer/sec:     29.02MB

    Running 30s test @ http://127.0.0.1:8080/
      4 threads and 16 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   220.12us   96.77us   1.98ms   64.61%
        Req/Sec    18.08k   234.07    18.71k    82.06%
      Latency Distribution
         50%  223.00us
         75%  295.00us
         90%  342.00us
         99%  397.00us
      2165220 requests in 30.10s, 0.85GB read

* 512 connections / 2 server threads / 4 wrk threads: ::

    Requests/sec:  74050.48
    Transfer/sec:     29.87MB

    Running 30s test @ http://127.0.0.1:8080/
      4 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     6.86ms    6.06ms 219.10ms   54.85%
        Req/Sec    18.64k     1.62k   36.19k    91.42%
      Latency Distribution
         50%    7.25ms
         75%   12.54ms
         90%   13.56ms
         99%   14.84ms
      2225585 requests in 30.05s, 0.88GB read
      Socket errors: connect 0, read 89, write 0, timeout 0

* 2048 connections / 2 server threads / 4 wrk threads: ::

    Requests/sec:  74714.23
    Transfer/sec:     30.14MB

    Running 30s test @ http://127.0.0.1:8080/
      4 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    52.45ms   87.02ms 997.26ms   88.24%
        Req/Sec    18.84k     3.18k   35.31k    80.77%
      Latency Distribution
         50%   23.60ms
         75%   34.86ms
         90%  162.92ms
         99%  435.41ms
      2244296 requests in 30.04s, 0.88GB read
      Socket errors: connect 0, read 106, write 0, timeout 51


Results notes
.............

* the machine was my personal laptop:  6 years old with an Intel Core i7 3667U (2 cores with 2 threads each);
* the ``kawipiko-server`` was started with ``--processes 1 --threads 2``;  (i.e. 2 threads handling the requests;)
* the ``kawipiko-server`` was started with ``--archive-inmem``;  (i.e. the CDB database file was preloaded into memory, thus no disk I/O;)
* the benchmarking tool was wrk_;
* both ``kawipiko-server`` and ``wrk`` tools were run on the same machine;
* both ``kawipiko-server`` and ``wrk`` tools were pinned on different physical cores;
* the benchmark was run over loopback networking (i.e. ``127.0.0.1``);
* the served file contains the content ``Hello World!``;
* the protocol was HTTP (i.e. no TLS), with keep-alive;
* see the `methodology section <#methodology>`__ for details;




Comparisons
-----------


Comparisons with NGinx
......................

* NGinx 512 connections / 2 server workers / 4 wrk thread: ::

    Requests/sec:  97910.36
    Transfer/sec:     24.56MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      4 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     5.11ms    1.30ms  17.59ms   85.08%
        Req/Sec    24.65k     1.35k   42.68k    78.83%
      Latency Distribution
         50%    5.02ms
         75%    5.32ms
         90%    6.08ms
         99%    9.62ms
      2944219 requests in 30.07s, 738.46MB read

* NGinx 2048 connections / 2 server workers / 4 wrk thread: ::

    Requests/sec:  93240.70
    Transfer/sec:     23.39MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      4 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    36.33ms   56.44ms 859.65ms   90.18%
        Req/Sec    23.61k     6.24k   51.88k    74.33%
      Latency Distribution
         50%   19.25ms
         75%   25.46ms
         90%   89.69ms
         99%  251.04ms
      2805639 requests in 30.09s, 703.70MB read
      Socket errors: connect 0, read 25, write 0, timeout 66

* (the NGinx configuration file can be found in the `examples folder <./examples>`__;  the configuration was obtained after many experiments to squeeze out of NGinx as much performance as possible, given the targeted use-case, namely many small static files;)


Comparisons with others
.......................

* darkhttpd_ 512 connections / 1 server process / 4 wrk threads: ::

    Requests/sec:  38191.65
    Transfer/sec:      8.74MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      4 threads and 512 connections
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
        --timeout 6s \
        --duration 30s \
        --latency \
        http://127.0.0.1:8080/ \
    #

* 4096 concurrent connections (handled by 4 threads): ::

    wrk \
        --threads 4 \
        --connections 4096 \
        --timeout 6s \
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
    prlimit -n16384 -- \
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
    prlimit -n16384 -- \
    \
    sudo -u "${USER}" -n -E -P -- \
    \
    wrk \
        ... \
    #




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

