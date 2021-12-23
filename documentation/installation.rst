

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




Installation
============

.. contents::
    :depth: 2
    :local:
    :backlinks: none




--------




Download prebuilt executables
-----------------------------


.. warning ::

  No executables are currently available for download!

  Please consult the `build from sources section <#build-from-sources>`__ for now.




--------




Build from sources
------------------




Install the prerequisites
.........................


* Ubuntu / Debian: ::

    apt-get install git-core
    apt-get install golang


* OpenSUSE: ::

    zypper install git-core
    zypper install go


* other Linux / FreeBSD / OpenBSD / OSX:

  * fetch and install Go from: `<https://golang.org/dl>`__
  * add ``/usr/local/go/bin`` to your ``PATH``;
  * install Git;




Prepare the environment
.......................


::

    mkdir -- \
            /tmp/kawipiko \
            /tmp/kawipiko/bin \
            /tmp/kawipiko/src \
            /tmp/kawipiko/go \
    #




Fetch the sources
.................


Either clone the full Git repository: ::

    git clone \
            -b development \
            git://github.com/volution/kawipiko.git \
            /tmp/kawipiko/src \
    #


Either fetch and extract the latest sources bundle: ::

    curl \
            -s -S -f \
            -o /tmp/kawipiko/src.tar.gz \
            https://codeload.github.com/volution/kawipiko/tar.gz/development \
    #

    tar \
            -x -z -v \
            -f /tmp/kawipiko/src.tar.gz \
            -C /tmp/kawipiko/src \
            --strip-components 1 \
    #




Build the dynamic executables
.............................


Compile the (dynamic) executables: ::

    cd /tmp/kawipiko/src/sources

    #### build `kawipiko` dynamic all-in-one executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko \
            ./cmd/wrapper.go \
    #

    #### build `kawipiko-server` dynamic executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko-server \
            ./cmd/server.go \
    #

    #### build `kawipiko-archiver` dynamic executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko-archiver \
            ./cmd/archiver.go \
    #




Build the static executables
............................


Compile the (static) executables: ::

    cd /tmp/kawipiko/src/sources

    #### build `kawipiko` static all-in-one executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -tags 'netgo' \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko \
            ./cmd/wrapper.go \
    #

    #### build `kawipiko-server` static executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -tags 'netgo' \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko-server \
            ./cmd/server.go \
    #

    #### build `kawipiko-archiver` static executable
    env \
            GOPATH=/tmp/kawipiko/go \
    go build \
            -tags 'netgo' \
            -gcflags 'all=-l=4' \
            -ldflags 'all=-s' \
            -trimpath \
            -o /tmp/kawipiko/bin/kawipiko-archiver \
            ./cmd/archiver.go \
    #




Deploy the executables
......................


Just copy the two executables anywhere on the system, or any compatible remote system: ::

    cp \
            -t /usr/local/bin \
            /tmp/kawipiko/bin/kawipiko-server \
            /tmp/kawipiko/bin/kawipiko-archiver \
    #

