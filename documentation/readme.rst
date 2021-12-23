

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




``kawipiko`` is a **lightweight static HTTP server** written in Go;
focused on serving static content **as fast and efficient as possible**,
with the **lowest latency**, and with the lowest resource consumption (either CPU, RAM, IO);
supporting both **HTTP/1 (with or without TLS), HTTP/2 and HTTP/3 (over QUIC)**;
available as a **single statically linked executable** without any other dependencies.


However, *simple* doesn't imply *dumb* or *limited*,
instead it implies *efficient* through the removal of superfluous features,
thus being inline with UNIX's old philosophy of
"`do one thing and do it well <https://en.wikipedia.org/wiki/Unix_philosophy#Do_One_Thing_and_Do_It_Well>`__".
Therefore, it supports only ``GET`` requests,
and does not provide features like dynamic content generation, authentication, reverse proxying, etc.;
meanwhile still providing compression (``gzip``, ``zopfli``, or ``brotli``),
plus HTML-CSS-JS minifying (TODO),
without affecting its performance
(due to its unique architecture as described below).


What ``kawipiko`` does provide is something very unique, that no other HTTP server offers:
the static content is served from a `CDB file <#why-cdb>`__ with almost no latency
(as compared to classical static servers that still have to pass through the OS via the ``open-read-close`` syscalls).
Moreover, as noted earlier, the static content can still be compressed or minified ahead of time,
thus reducing not only CPU but also bandwidth and latency.


`CDB files <#why-cdb>`__ are binary database files that provide efficient read-only key-value lookup tables,
initially used in some DNS and SMTP servers,
mainly for their low overhead lookup operations,
zero locking in multi-threaded / multi-process scenarios,
and "atomic" multi-record updates.
This also makes them suitable for low-latency static content serving over HTTP,
which is what this project provides.


For those familiar with Netlify (or competitors like CloudFlare Pages, GitHub Pages, etc.),
``kawipiko`` is a *host-it-yourself* alternative featuring:

* self-contained deployment with simple configuration;
  (i.e. just `fetch the executable <#installation>`__ and use the `proper flags <#kawipiko-server>`__;)

* low and constant resource consumption (both in terms of CPU and RAM);
  (i.e. you won't have surprises when under load;)

* (hopefully) extremely secure;
  (i.e. it doesn't launch processes, it doesn't connect to other services or databases, it doesn't open any files, etc.;
  basically you can easily ``chroot`` it, or containerize it as is in fashion these days;)

* highly portable, supporting at least Linux (the main development, testing and deployment platform), FreeBSD, OpenBSD, and OSX;


For a complete list of features please consult the `features section <#features>`__.
Unfortunately, there are also some tradeoffs as described in the `limitations section <#limitations>`__
(although none are critical).


With regard to performance, as described in the `benchmarks section <#benchmarks>`__,
``kawipiko`` is at least on-par with NGinx,
sustaining over 100K requests / second with 0.25ms latency for 99% of the requests even on my 6 years old laptop.
However the main advantage over NGinx is not raw performance,
but deployment and configuration simplicity,
plus efficient management and storage of large collections of many small files.




--------




.. contents::
    :depth: 1
    :backlinks: none




--------




.. image:: ./documentation/banner.png
     :width: 50em




--------




Manual
======

.. contents::
    :local:
    :backlinks: none




Workflow
--------


The project provides the following executables (statically linked, without any other dependencies):

* ``kawipiko-server`` -- which serves the static content from the CDB archive either via HTTP (with or without TLS), HTTP/2 or HTTP/3 (over QUIC);

* ``kawipiko-archiver`` -- which creates the CDB archive from a source folder holding the static content,
  optionally compressing and minifying files;

* ``kawipiko`` -- an all-in-one executable that bundles all functionality in one executable;
  (i.e. ``kawipiko server ...`` or ``kawipiko archiver ...``);


Unlike most (if not all) other servers out-there,
in which you just point your web server to the folder holding the static website content root,
``kawipiko`` takes a radically different approach:
in order to serve the static content,
one has to first *archive* the content into the CDB archive through ``kawipiko-archiver``,
and then one can *serve* it from the CDB archive through ``kawipiko-server``.


This two step phase also presents a few opportunities:

* one can decouple the "building", "testing", and "publishing" phases of a static website,
  by using a similar CI/CD pipeline as done for other software projects;

* one can instantaneously rollback to a previous version if the newly published one has issues;

* one can apply extreme compression (e.g. ``zopfli`` or ``brotli``),
  to trade CPU during deployment vs latency and bandwidth at runtime.




kawipiko-server
---------------


See the `dedicated manual <./documentation/manuals/server.rst>`__.

This document is also available
in `plain text <./documentation/manuals/server.txt>`__,
or as a `man page <./documentation/manuals/server.1.man>`__.




kawipiko-archiver
-----------------


See the `dedicated manual <./documentation/manuals/archiver.rst>`__.

This document is also available
in `plain text <./documentation/manuals/archiver.txt>`__,
or as a `man page <./documentation/manuals/archiver.1.man>`__.




--------




Examples
========


* fetch and extract the Python 3.10 documentation HTML archive: ::

    curl \
            -s -S -f \
            -o ./python-3.10.1-docs-html.tar.bz2 \
            https://docs.python.org/3/archives/python-3.10.1-docs-html.tar.bz2 \
    #

    tar \
            -x -j -v \
            -f ./python-3.10.1-docs-html.tar.bz2 \
    #


* create the CDB archive (without any compression): ::

    kawipiko-archiver \
            --archive ./python-3.10.1-docs-html-nocomp.cdb \
            --sources ./python-3.10.1-docs-html \
            --debug \
    #


* create the CDB archive (with ``gzip`` compression): ::

    kawipiko-archiver \
            --archive ./python-3.10.1-docs-html-gzip.cdb \
            --sources ./python-3.10.1-docs-html \
            --compress gzip \
            --debug \
    #


* create the CDB archive (with ``zopfli`` compression): ::

    kawipiko-archiver \
            --archive ./python-3.10.1-docs-html-zopfli.cdb \
            --sources ./python-3.10.1-docs-html \
            --compress zopfli \
            --debug \
    #


* create the CDB archive (with ``brotli`` compression): ::

    kawipiko-archiver \
            --archive ./python-3.10.1-docs-html-brotli.cdb \
            --sources ./python-3.10.1-docs-html \
            --compress brotli \
            --debug \
    #


* serve the CDB archive (with ``gzip`` compression): ::

    kawipiko-server \
            --bind 127.0.0.1:8080 \
            --archive ./python-3.10.1-docs-html-gzip.cdb \
            --archive-mmap \
            --archive-preload \
            --debug \
    #


* compare sources and archive sizes: ::

    du \
            -h -s \
            \
            ./python-3.10.1-docs-html-nocomp.cdb \
            ./python-3.10.1-docs-html-gzip.cdb \
            ./python-3.10.1-docs-html-zopfli.cdb \
            ./python-3.10.1-docs-html-brotli.cdb \
            \
            ./python-3.10.1-docs-html \
            ./python-3.10.1-docs-html.tar.bz2 \
    #

    45M     ./python-3.10.1-docs-html-nocomp.cdb
    9.7M    ./python-3.10.1-docs-html-gzip.cdb
    ???     ./python-3.10.1-docs-html-zopfli.cdb
    7.9M    ./python-3.10.1-docs-html-brotli.cdb

    46M     ./python-3.10.1-docs-html
    6.0M    ./python-3.10.1-docs-html.tar.bz2




--------




Installation
============


See the `dedicated installation document <./documentation/installation.rst>`__.




--------




Features
========

.. contents::
    :local:
    :backlinks: none




Implemented
-----------


The following is a list of the most important features:

* (optionally)  the static content is compressed or minified when the CDB archive is created,
  thus no CPU cycles are used while serving requests;

* (optionally)  the static content can be compressed with either ``gzip``, ``zopfli`` or ``brotli``;

* (optionally)  in order to reduce the serving latency even further,
  one can preload the entire CDB archive in memory, or alternatively mapping it in memory (using ``mmap``);
  this trades memory for CPU;

* (optionally)  caching the static content fingerprint and compression,
  thus significantly reducing the CDB archive rebuilding time,
  and significantly reducing the IO for the source file-system;

* atomic static website content changes;
  because the entire content is held in a single CDB archive,
  and because the file replacement is atomically achieved via the ``rename`` syscall (or the ``mv`` tool),
  all served resources are observed to change at the same time;

* ``_wildcard.*`` files (where ``.*`` are the regular extensions like ``.txt``, ``.html``, etc.)
  which will be used if an actual resource is not found under that folder;
  (these files respect the hierarchical tree structure, i.e. "deeper" ones override the ones closer to "root";)

* support for HTTP/1 (with or without TLS), by leveraging ``github.com/valyala/fasthttp``;

* support for HTTP/2, by leveraging Go's ``net/http``;

* support for HTTP/3 (over QUIC), by leveraging ``github.com/lucas-clemente/quic-go``;




Pending
-------


The following is a list of the most important features that are currently missing and are planed to be implemented:

* (TODO)  support for custom HTTP response headers (for specific files, for specific folders, etc.);
  (currently only ``Content-Type``, ``Content-Length``, ``Content-Encoding`` are included;
  additionally ``Cache-Control: public, immutable, max-age=3600``, optionally ``ETag``,
  and a few TLS or security related headers can also be included;)

* (TODO)  support for mapping virtual hosts to key prefixes;
  (currently virtual hosts, i.e. the ``Host`` header, are ignored;)

* (TODO)  support for mapping virtual hosts to multiple CDB archives;
  (i.e. the ability to serve multiple domains, each with its own CDB archive;)

* (TODO)  automatic reloading of the CDB archives;

* (TODO)  minifying HTML, CSS and JavaScript, by leveraging ``https://github.com/tdewolff/minify``;

* (TODO)  customized error pages (embedded in the CDB archive);




Limitations
-----------


As stated in the `about section <#about>`__, nothing comes for free,
and in order to provide all these features, some corners had to be cut:

* (TODO)  currently if the CDB archive changes,
  the server needs to be restarted in order to pickup the changed files;

* (won't fix)  the CDB archive **maximum size is 4 GiB** (after compression and minifying),
  and there can't be more than 16M resources;
  (however if you have a static website this large,
  you are probably doing something extremely wrong,
  as large files should be offloaded to something like AWS S3,
  and served through a CDN like CloudFlare or AWS CloudFront;)

* (won't fix)  the server **does not support per-request decompression / recompression**;
  this implies that if the content was saved in the CDB archive with compression (say ``brotli``),
  the server will serve all resources compressed (i.e. ``Content-Encoding: brotli``),
  regardless of what the browser accepts (i.e. ``Accept-Encoding: gzip``);
  the same applies for uncompressed content;
  (however always using ``gzip`` compression is safe enough,
  as it is implemented in virtually all browsers and HTTP clients out there;)

* (won't fix)  regarding the "atomic" static website changes,
  there is a small time window in which a client that has fetched an "old" version of a resource (say an HTML page),
  but it has not yet fetched the required resources (say the CSS or JS files),
  and in between fetching the HTML and CSS/JS the CDB archive was changed,
  the client will consequently fetch the new version of these required resources;
  however due to the low latency serving, this time window is extremely small;
  (**this is not a limitation of this HTTP server, but a limitation of the way websites are built;**
  always use fingerprints in your resources URL,
  and perhaps always include the current and previous version on each deploy;)




--------




Benchmarks
==========


See the `dedicated benchmarks document <./documentation/benchmarks.rst>`__.




--------




FAQ
===




Is it production ready?
-----------------------


Yes, it currently is serving ~600K HTML pages.


Although, being open source, you are responsible for making sure it works within your requirements!


However, I am available for consulting on its deployment and usage.  :)




Why CDB?
--------


Until I expand upon why I have chosen to use CDB for service static website content,
you can read about the `sparkey <https://github.com/spotify/sparkey>`__ from Spotify.




Why Go?
-------


Because Go is highly portable, highly stable,
and especially because it can easily support cross-compiling statically linked binaries
to any platform it supports.




Why not Rust?
-------------


Because Rust fails to easily support cross-compiling (statically or dynamically linked) executables
to any platform it supports.


Because Rust is less portable than Go;
for example Rust doesn't consider OpenBSD as a "tier-1" platform.




--------




Notice (copyright and licensing)
================================

.. contents::
    :local:
    :backlinks: none




Authors
-------


Ciprian Dorin Craciun
  * `ciprian@volution.ro <mailto:ciprian@volution.ro>`__
    or `ciprian.craciun@gmail.com <mailto:ciprian.craciun@gmail.com>`__
  * `<https://volution.ro/ciprian>`__
  * `<https://github.com/cipriancraciun>`__




Notice -- short version
-----------------------


The code is licensed under AGPL 3 or later.


If you **change** the code within this repository **and use** it for **non-personal** purposes,
you'll have to release it as per AGPL.




Notice -- long version
----------------------


For details about the copyright and licensing,
please consult the `notice <./documentation/licensing/notice.txt>`__ file
in the `documentation/licensing <./documentation/licensing>`__ folder.


If someone requires the sources and/or documentation to be released
under a different license, please send an email to the authors,
stating the licensing requirements, accompanied with the reasons
and other details; then, depending on the situation, the authors might
release the sources and/or documentation under a different license.




--------




References
==========


See the `dedicated references document <./documentation/references.rst>`__.

