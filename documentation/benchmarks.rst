

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




Benchmarks
==========

.. contents::
    :depth: 2
    :local:
    :backlinks: none




--------




Summary
-------


.. important::

  Last updated in December 2021!


Bottom line (**even on my 6 years old laptop**),
using only 1 core with 2 hyperthreads
(one core for the server, and a separate core for the load generator),
with HTTP Keep-Alive capped at 256k requests per connection:

* under normal conditions (16 concurrent connections),
  I get around 105k requests / second,
  at about 0.4ms latency for 99% of the requests;

* under normal conditions (64 concurrent connections),
  I get around 107k requests / second,
  at about 1.5ms latency for 99% of the requests;

* under light stress conditions (128 concurrent connections),
  I get around 110k requests / second,
  at about 3.0ms latency for 99% of the requests;

* under medium stress conditions (512 concurrent connections),
  I get around 104k requests / second,
  at about 9.3ms latency for 99% of the requests
  (meanwhile the average is under 5.0ms);

* under high stress conditions (2048 concurrent connections),
  I get around 103k requests / second,
  at about 240ms latency for 99% of the requests
  (meanwhile the average is under 20ms);

* under extreme stress conditions (16384 concurrent connections)
  (i.e. someone tries to DDOS the server),
  I get around 90k requests / second,
  at about 3.1s latency for 99% of the requests
  (meanwhile the average is under 200ms);

* **the performance is at least on-par with NGinx**;
  however, especially for a real world scenarios
  (i.e. thousand of small files, accessed in a random patterns),
  I believe ``kawipiko`` fares much better;
  (not to mention how simple it is to configure and deploy ``kawipiko`` as compared to NGinx,
  which took a lot of time, fiddling, and trial and error to get it right;)


Regarding HTTPS, my initial benchmarks
(only covering plain HTTPS with HTTP/1)
seem to indicate that ``kawipiko`` is at least on-par with NGinx.


Regarding HTTP/2, my initial benchmarks
seem to indicate that ``kawipiko``'s performance is 6 times less than plain HTTPS with HTTP/1
(mainly due to the unoptimized Go ``net/http`` implementation).
In this regard NGinx is much better, having a HTTP/2 performance similar to plain HTTPS with HTTP/1.


Regarding HTTP/3, given that the QUIC library is still experimental,
my initial benchmarks seem to indicate that ``kawipiko``'s performance is quite poor
(at about 5k requests / second).




--------




Performance
-----------


.. important::

  Last updated in August 2018!

  The results are based on an older version of ``kawipiko``;
  the current version is at least 10% more efficient.

  The methodology used is described in a `dedicated section <#methodology>`__.




Performance values
..................


.. note ::

  Please note that the values under *Thread Stats* are reported per thread.
  Therefore it is best to look at the first two values, i.e. *Requests/sec*.


* ``kawipiko``, 16 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec: 111720.73
    Transfer/sec:     18.01MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 16 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   139.36us   60.27us   1.88ms   64.91%
        Req/Sec    56.14k   713.04    57.60k    91.36%
      Latency Distribution
         50%  143.00us      75%  184.00us
         90%  212.00us      99%  261.00us
      3362742 requests in 30.10s, 541.98MB read


* ``kawipiko``, 128 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec: 118811.41
    Transfer/sec:     19.15MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 128 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     1.03ms  705.69us  19.53ms   63.54%
        Req/Sec    59.71k     1.69k   61.70k    96.67%
      Latency Distribution
         50%    0.99ms      75%    1.58ms
         90%    1.89ms      99%    2.42ms
      3564527 requests in 30.00s, 574.50MB read


* ``kawipiko``, 512 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec: 106698.89
    Transfer/sec:     17.20MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     4.73ms    3.89ms  39.32ms   39.74%
        Req/Sec    53.71k     1.73k   69.18k    84.33%
      Latency Distribution
         50%    4.96ms      75%    8.63ms
         90%    9.19ms      99%   10.30ms
      3206540 requests in 30.05s, 516.80MB read
      Socket errors: connect 0, read 105, write 0, timeout 0


* ``kawipiko``, 2048 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec: 100296.65
    Transfer/sec:     16.16MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    45.42ms   85.14ms 987.70ms   88.62%
        Req/Sec    50.61k     5.59k   70.14k    71.74%
      Latency Distribution
         50%   16.30ms      75%   28.44ms
         90%  147.60ms      99%  417.40ms
      3015868 requests in 30.07s, 486.07MB read
      Socket errors: connect 0, read 128, write 0, timeout 86


* ``kawipiko``, 4096 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec:  95628.34
    Transfer/sec:     15.41MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 4096 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    90.50ms  146.08ms 999.65ms   88.49%
        Req/Sec    48.27k     6.09k   66.05k    76.34%
      Latency Distribution
         50%   23.31ms      75%  112.06ms
         90%  249.41ms      99%  745.94ms
      2871404 requests in 30.03s, 462.79MB read
      Socket errors: connect 0, read 27, write 0, timeout 4449


* ``kawipiko``, 16384 connections / 2 server threads / 2 ``wrk`` threads: ::

    Requests/sec:  53548.52
    Transfer/sec:      8.63MB

    Running 30s test @ http://127.0.0.1:8080/
      2 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   206.21ms  513.75ms   6.00s    92.56%
        Req/Sec    31.37k     5.68k   44.44k    76.13%
      Latency Distribution
         50%   35.38ms      75%   62.78ms
         90%  551.33ms      99%    2.82s
      1611294 requests in 30.09s, 259.69MB read
      Socket errors: connect 0, read 115, write 0, timeout 2288




Performance notes
.................


* the machine was my personal laptop,
  with an Intel Core i7 3667U (2 physical cores times 2 hyper-threads each),
  and with 8 GiB of RAM;

* the ``kawipiko-server`` was started with ``--processes 1 --threads 2``;
  (i.e. 2 threads handling the requests;)

* the ``kawipiko-server`` was started with ``--archive-inmem``;
  (i.e. the CDB database file was preloaded into memory, thus no disk IO;)

* the ``kawipiko-server`` was started with ``--security-headers-disable``;
  (because these headers are not set by default by other HTTP servers;)

* the ``kawipiko-server`` was started with ``--timeout-disable``;
  (because, due to a known Go issue, using ``net.Conn.SetDeadline`` has an impact of about 20% of the raw performance;
  thus the reported values above might be about 10%-15% smaller when used with timeouts;)

* the benchmarking tool was ``wrk``;

* both ``kawipiko-server`` and ``wrk`` tools were run on the same machine;

* both ``kawipiko-server`` and ``wrk`` tools were pinned on different physical cores;

* the benchmark was run over loopback networking (i.e. ``127.0.0.1``);

* the served file contains ``Hello World!``;

* the protocol was HTTP (i.e. no TLS), with keep-alive;

* both the CDB and the NGinx folder were put on ``tmpfs``
  (which implies that the disk is not a limiting factor);
  (in fact ``kawipiko`` performs quite well even on spinning disks due to careful storage management;)




--------




Comparisons
-----------


.. important::

  Last updated in August 2019!

  The results are based on an older version of ``kawipiko``;
  the current version is at least 10% more efficient.

  The methodology used is described in a `dedicated section <#methodology>`__.




Comparisons with NGinx
......................


* NGinx, 512 connections / 2 worker processes / 2 ``wrk`` threads: ::

    Requests/sec:  79816.08
    Transfer/sec:     20.02MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     6.07ms    1.90ms  19.83ms   71.67%
        Req/Sec    40.17k     1.16k   43.35k    69.83%
      Latency Distribution
         50%    6.13ms      75%    6.99ms
         90%    8.51ms      99%   11.10ms
      2399069 requests in 30.06s, 601.73MB read


* NGinx, 2048 connections / 2 worker processes / 2 ``wrk`` threads: ::

    Requests/sec:  78211.46
    Transfer/sec:     19.62MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 2048 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    27.11ms   20.27ms 490.12ms   97.76%
        Req/Sec    39.45k     2.45k   49.98k    70.74%
      Latency Distribution
         50%   24.80ms      75%   29.67ms
         90%   34.99ms      99%  126.97ms
      2351933 requests in 30.07s, 589.90MB read
      Socket errors: connect 0, read 0, write 0, timeout 11


* NGinx, 4096 connections / 2 worker processes / 2 ``wrk`` threads: ::

    Requests/sec:  75970.82
    Transfer/sec:     19.05MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 4096 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    70.25ms   73.68ms 943.82ms   87.21%
        Req/Sec    38.37k     3.79k   49.06k    70.30%
      Latency Distribution
         50%   46.37ms      75%   58.28ms
         90%  179.08ms      99%  339.05ms
      2282223 requests in 30.04s, 572.42MB read
      Socket errors: connect 0, read 0, write 0, timeout 187


* NGinx, 16384 connections / 2 worker processes / 2 ``wrk`` threads: ::

    Requests/sec:  43909.67
    Transfer/sec:     11.01MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   223.87ms  551.14ms   5.94s    92.92%
        Req/Sec    32.95k    13.35k   51.56k    76.71%
      Latency Distribution
         50%   32.62ms      75%  222.93ms
         90%  558.04ms      99%    3.17s
      1320562 requests in 30.07s, 331.22MB read
      Socket errors: connect 0, read 12596, write 34, timeout 1121


* the NGinx configuration file can be found in the `examples folder <../examples/nginx>`__;
  the configuration was obtained after many experiments to squeeze out of NGinx as much performance as possible,
  given the targeted use-case, namely many small files;

* moreover NGinx seems to be quite sensitive to the actual path requested:

    * if one requests ``http://127.0.0.1:8080/``,
      and one has configured NGinx to look for ``index.txt``,
      and that file actually exists,
      the performance is quite a bit lower than just asking for that file;
      (perhaps it issues more syscalls searching for the index file;)

    * if one requests ``http://127.0.0.1:8080/index.txt``,
      as mentioned above, it achieves the higher performance;
      (perhaps it issues fewer syscalls;)

    * if one requests ``http://127.0.0.1:8080/does-not-exist``,
      it seems to achieve the best performance;
      (perhaps it issues the least amount of syscalls;)
      (however this is not an actual useful corner-case;)

    * it must be noted that ``kawipiko`` doesn't exhibit this behaviour,
      the same performance is achieved regardless of the path variant;

    * therefore the benchmarks above use ``/index.txt`` as opposed to ``/``,
      in order not to disfavour NGinx;




Comparisons with others
.......................


* ``darkhttpd``, 512 connections / 1 server process / 2 ``wrk`` threads: ::

    Requests/sec:  38191.65
    Transfer/sec:      8.74MB

    Running 30s test @ http://127.0.0.1:8080/index.txt
      2 threads and 512 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    17.51ms   17.30ms 223.22ms   78.55%
        Req/Sec     9.62k     1.94k   17.01k    72.98%
      Latency Distribution
         50%    7.51ms      75%   32.51ms
         90%   45.69ms      99%   53.00ms
      1148067 requests in 30.06s, 262.85MB read




--------




OpenStreetMap tiles
-------------------


.. important::

  Last updated in August 2019!

  The results are based on an older version of ``kawipiko``;
  the current version is at least 10% more efficient.

  The methodology used is described in a `dedicated section <#methodology>`__.




Scenario notes
..............


As a benchmark much closer to the "real world" use-cases for ``kawipiko`` I've done the following:

* downloaded from OpenStreetMap servers all tiles for my home town
  (from zoom level 0 to zoom level 19), which resulted in:

  * around ~250k PNG files totaling ~330 MiB;
  * with an average of 1.3 KiB and a median of 103B;
    (i.e. lots of extreemly small files;)
  * occupying actualy around 1.1 GiB of storage (on Ext4) due to file-system overheads;

* created a CDB archive, which resulted in:

  * a single file totaling ~376 MiB (both "apparent" and "occupied" storage);
    (i.e. no storage space wasted;)
  * which contains only ~100k PNG files, due to elimination of duplicate PNG files;
    (i.e. at higher zoom levels, the tiles start to repeat;)

* listed all the available tiles, and benchmarked both ``kawipiko`` and NGinx,
  with 16k concurrent connections;

* the methodology is the same one described above,
  with the following changes:

  * the machine used was my desktop,
    with an Intel Core i7 4770 (4 physical cores times 2 hyper-threads each),
    and with 32 GiB of RAM;

  * the files (both CDB and tiles folder) were put in ``tmpfs``;

  * both ``kawipiko``, NGinx and ``wrk``
    were configured to use 8 threads,
    and were pinned on two separate physical cores each;

  * (the machine had almost nothing running on it except the minimal required services;)




Results notes
.............


Based on my benchmark the following are my findings:

* ``kawipiko`` outperformed NGinx by ~25% in requests / second;

* ``kawipiko`` outperformed NGinx by ~29% in average response latency;

* ``kawipiko`` outperformed NGinx by ~40% in 90-percentile response latency;

* ``kawipiko`` used ~6% less CPU while serving requests for 2 minutes;

* ``kawipiko`` used ~25% less CPU per request;

* NGinx used the least amount of RAM,
  meanwhile ``kawipiko`` used around 1 GiB of RAM
  (due to either in RAM loading or ``mmap`` usage);




Results values
..............


* ``kawipiko`` with ``--archive-inmem`` and ``--index-all`` (1 process, 8 threads): ::

    Requests/sec: 238499.86
    Transfer/sec:    383.59MB

    Running 2m test @ http://127.9.185.194:8080/
      8 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   195.39ms  412.84ms   5.99s    92.33%
        Req/Sec    30.65k    10.20k  213.08k    79.41%
      Latency Distribution
         50%   28.02ms      75%  221.17ms
         90%  472.41ms      99%    2.19s
      28640139 requests in 2.00m, 44.98GB read
      Socket errors: connect 0, read 0, write 0, timeout 7032


* ``kawipiko`` with ``--archive-mmap`` (1 process, 8 threads): ::

    Requests/sec: 237239.35
    Transfer/sec:    381.72MB

    Running 2m test @ http://127.9.185.194:8080/
      8 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   210.44ms  467.84ms   6.00s    92.57%
        Req/Sec    30.77k    12.29k  210.17k    86.67%
      Latency Distribution
         50%   26.51ms      75%  221.63ms
         90%  494.93ms      99%    2.67s
      28489533 requests in 2.00m, 44.77GB read
      Socket errors: connect 0, read 0, write 0, timeout 10730


* ``kawipiko`` with ``--archive-mmap`` (8 processes, 1 thread): ::

    Requests/sec: 248266.83
    Transfer/sec:    399.29MB

    Running 2m test @ http://127.9.185.194:8080/
      8 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   209.30ms  469.05ms   5.98s    92.25%
        Req/Sec    31.86k     8.58k   83.99k    69.93%
      Latency Distribution
         50%   23.08ms      75%  215.28ms
         90%  502.80ms      99%    2.64s
      29816650 requests in 2.00m, 46.83GB read
      Socket errors: connect 0, read 0, write 0, timeout 15244


* NGinx (8 worker processes): ::

    Requests/sec: 188255.32
    Transfer/sec:    302.88MB

    Running 2m test @ http://127.9.185.194:8080/
      8 threads and 16384 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   266.18ms  538.72ms   5.93s    90.78%
        Req/Sec    24.15k     8.34k  106.48k    74.56%
      Latency Distribution
         50%   34.34ms      75%  253.57ms
         90%  750.29ms      99%    2.97s
      22607727 requests in 2.00m, 35.52GB read
      Socket errors: connect 0, read 109, write 0, timeout 16833




--------




Methodology
-----------


* get the ``kawipiko`` executables (either `download <./installation.rst#download-prebuilt-executables>`__ or `build <./installation.rst#build-from-sources>`__ them);

* get the ``hello-world.cdb`` (from the `examples <../examples>`__ folder inside the repository);

* install NGinx and ``wrk`` from the distribution packages;




Single process / single threaded
................................


* this scenario will yield a base-line performance per core;


* execute the server (in-memory and indexed)
  (i.e. the *best case scenario*): ::

    kawipiko-server \
            --bind 127.0.0.1:8080 \
            --archive ./hello-world.cdb \
            --archive-inmem \
            --index-all \
            --processes 1 \
            --threads 1 \
    #


* execute the server (memory mapped)
  (i.e. the *the recommended scenario*): ::

    kawipiko-server \
            --bind 127.0.0.1:8080 \
            --archive ./hello-world.cdb \
            --archive-mmap \
            --processes 1 \
            --threads 1 \
    #




Single process / two threads
............................


* this scenario is the usual setup;
  configure ``--threads`` to equal the number of logical cores
  (i.e. multiply the number of physical cores with
  the number of hyper-threads per physical core);


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


* ``wrk``, 512 concurrent connections, handled by 2 threads: ::

    wrk \
            --threads 2 \
            --connections 512 \
            --timeout 1s \
            --duration 30s \
            --latency \
            http://127.0.0.1:8080/index.txt \
    #


* ``wrk``, 4096 concurrent connections, handled by 2 threads: ::

    wrk \
            --threads 2 \
            --connections 4096 \
            --timeout 1s \
            --duration 30s \
            --latency \
            http://127.0.0.1:8080/index.txt \
    #




Methodology notes
.................


* the number of threads for the server plus for ``wrk`` shouldn't be larger than the number of available cores;
  (or use different machines for the server and the client;)

* also take into account that by default the number of file descriptors
  on most UNIX / Linux systems is 1024,
  therefore if you want to try with more connections than 1000,
  you need to raise this limit;
  (see bellow;)

* additionally, you can try to pin the server and ``wrk`` to specific cores,
  increase various priorities (scheduling, IO, etc.);
  (given that Intel processors have hyper-threading which appear to the OS as individual cores, you should make sure that you pin each process on cores part of the same physical processor / core;)

* pinning the server (cores ``0`` and ``1`` are mapped on the physical core ``1``): ::

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

* pinning the client (cores ``2`` and ``3`` are mapped on the physical core ``2``): ::

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

