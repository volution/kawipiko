

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


See the `releases page on GitHub <https://github.com/volution/kawipiko/releases>`__.

As a shortcut, the following are the self-contained and statically linked
all-in-one server and archiver executables
(for x86_64 / amd64 processors):

 * `<https://github.com/volution/kawipiko/releases/download/preview/kawipiko-wrapper--linux--v0.1.0--preview>`__
 * `<https://github.com/volution/kawipiko/releases/download/preview/kawipiko-wrapper--freebsd--v0.1.0--preview>`__
 * `<https://github.com/volution/kawipiko/releases/download/preview/kawipiko-wrapper--openbsd--v0.1.0--preview>`__
 * `<https://github.com/volution/kawipiko/releases/download/preview/kawipiko-wrapper--darwin--v0.1.0--preview>`__


For example, assuming one wants the ``preview`` version,
one can run the following commands: ::

    curl \
            -s -S -f -L \
            -o /tmp/kawipiko-server \
            https://github.com/volution/kawipiko/releases/download/preview/kawipiko-server--linux--v0.1.0--preview \
    #

    curl \
            -s -S -f -L \
            -o /tmp/kawipiko-archiver \
            https://github.com/volution/kawipiko/releases/download/preview/kawipiko-archiver--linux--v0.1.0--preview \
    #

    curl \
            -s -S -f -L \
            -o /tmp/kawipiko \
            https://github.com/volution/kawipiko/releases/download/preview/kawipiko-wrapper--linux--v0.1.0--preview \
    #

    chmod a=rx /tmp/kawipiko-server
    chmod a=rx /tmp/kawipiko-archiver
    chmod a=rx /tmp/kawipiko


One can replace ``preview`` with ``v0.x.y`` (see the releases page).


One can replace ``linux`` with ``freebsd``, ``openbsd`` or ``darwin``.




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

  * fetch and install Go from `<https://golang.org/dl>`__;
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
            https://github.com/volution/kawipiko.git \
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

    #### build `kawipiko` all-in-one dynamic executable
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

    #### build `kawipiko` all-in-one static executable
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

