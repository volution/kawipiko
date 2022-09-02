

#############################################
kawipiko -- blazingly fast static HTTP server
#############################################




``kawipiko-archiver``
---------------------


::

    >> kawipiko-archiver --help
    >> kawipiko-archiver --man

::

    --sources <path>

    --archive <path>

    --compress <gzip | zopfli | brotli | identity>
    --compress-level <number>
    --compress-cache <path>

    --exclude-index
    --exclude-strip
    --exclude-cache
    --include-etag

    --exclude-slash-redirects
    --include-folder-listing
    --exclude-paths-index

    --progress  --debug

    --version
    --help          (show this short help)
    --man           (show the full manual)

    --sources-md5   (dump an ``md5sum`` of the sources)
    --sources-cpio  (dump a ``cpio.gz`` of the sources)




--------




Flags
.....

``--sources``

    The path to the source folder that is the root of the static website content.

``--archive``

    The path to the target CDB file that contains the archived static content.

``--compress``, and ``--compress-level``

    Each individual file (and consequently of the corresponding HTTP response body) is compressed with either ``gzip``, ``zopfli`` or ``brotli``;  by default (or alternatively with ``identity``) no compression is used.

    Even if compression is explicitly requested, if the compression ratio is bellow a certain threshold (depending on the uncompressed size), the file is stored without any compression.
    (It's senseless to force the client to spend time and decompress the response body if that time is not recovered during network transmission.)

    The compression level can be chosen, the value depending on the algorithm:

    * ``gzip`` -- ``-1`` for algorithm default, ``-2`` for Huffman only, ``0`` to ``9`` for fast to slow;
    * ``zopfli`` -- ``-1`` for algorithm default, ``0`` to ``30`` iterations for fast to slow;
    * ``brotli`` -- ``-1`` for algorithm default, ``0`` to ``9`` for fast to slow, ``-2`` for extreme;
    * (by "algorithm default", it is meant "what that algorithm considers the recommended default compression level";)
    * ``kawipiko`` by default uses the maximum compression level for each algorithm;  (i.e. ``9`` for ``gzip``, ``30`` for ``zopfli``, and ``-2`` for ``brotli``;)

``--sources-cache <path>``, and ``--compress-cache <path>``

    At the given path a single file is created (that is an BBolt database), that will be used to cache the following information:

    * in case of ``--sources-cache``, the fingerprint of each file contents is stored, so that if the file was not changed, re-reading it shouldn't be attempted unless it is absolutely necessary;  also if the file is small enough, its contents is stored in this database (deduplicated by its fingerprint);
    * in case of ``--compress-cache`` the compression outcome of each file contents is stored (deduplicated by its fingerprint), so that compression is done only once over multiple runs;

    Each of these caches can be safely reused between multiple related archives, especially when they have many files in common.
    Each of these caches can be independently used (or shared).

    Using these caches allows one to very quickly rebuild an archive when only a couple of files have been changed, without even touching the file-system for the unchanged ones.

``--exclude-index``

    Disables using ``_index.*`` and ``index.*`` files (where ``.*`` is one of ``.html``, ``.htm``, ``.xhtml``, ``.xht``, ``.txt``, ``.json``, and ``.xml``) to respond to a request whose URL path ends in ``/`` (corresponding to the folder wherein ``_index.*`` or ``index.*`` file is located).
    (This can be used to implement "slash" blog style URL's like ``/blog/whatever/`` which maps to ``/blog/whatever/index.html``.)

``--exclude-strip``

    Disables using a file with the suffix ``.html``, ``.htm``, ``.xhtml``, ``.xht``, and ``.txt`` to respond to a request whose URL does not exactly match an existing file.
    (This can be used to implement "suffix-less" blog style URL's like ``/blog/whatever`` which maps to ``/blog/whatever.html``.)

``--exclude-cache``

    Disables adding an ``Cache-Control: public, immutable, max-age=3600`` header that forces the browser (and other intermediary proxies) to cache the response for an hour (the ``public`` and ``max-age=3600`` arguments), and furthermore not request it even on reloads (the ``immutable`` argument).

``--include-etag``

    Enables adding an ``ETag`` response header that contains the SHA256 of the response body.

    By not including the ``ETag`` header (i.e. the default), and because identical headers are stored only one, if one has many files of the same type (that in turn without ``ETag`` generates the same headers), this can lead to significant reduction in stored headers blocks, including reducing RAM usage.
    (At this moment it does not support HTTP conditional requests, i.e. the ``If-None-Match``, ``If-Modified-Since`` and their counterparts;  however this ``ETag`` header might be used in conjuction with ``HEAD`` requests to see if the resource has changed.)

``--exclude-slash-redirects``

    Disables adding redirects to/from paths with/without `/`
    (For example, by default, if `/file` exists, then there is also a `/file/` redirect towards `/file`;  and vice-versa from `/folder` towards `/folder/`.)

``--include-folder-listing``

    Enables the creation of an internal list of folders.

``--exclude-paths-index``

    Disables the creation of an internal list of references that can be used in conjunction with the ``--index-all`` flag of the ``kawipiko-server``.

``--progress``

    Enables periodic reporting of various metrics.

``--debug``

    Enables verbose logging.
    It will log various information about the archived files (including compression statistics).




Ignored files
.............

* any file with the following prefixes: ``.``, ``#``;
* any file with the following suffixes: ``~``, ``#``, ``.log``, ``.tmp``, ``.temp``, ``.lock``;
* any file that contains the following: ``#``;
* any file that exactly matches the following: ``Thumbs.db``, ``.DS_Store``;
* (at the moment these rules are not configurable through flags;)




Wildcard files
..............


By placing a file whose name matches ``_wildcard.*`` (i.e. with the prefix ``_wildcard.`` and any other suffix), it will be used to respond to any request whose URL fails to find a "better" match.

These wildcard files respect the folder hierarchy, in that wildcard files in (direct or transitive) subfolders override the wildcard file in their parents (direct or transitive).

In addition to ``_wildcard.*``, there is also support for ``_200.html`` (or just ``200.html``), plus ``_404.html`` (or just ``404.html``).




Redirect files
..............

By placing a file whose name is ``_redirects`` (or ``_redirects.txt``), it instructs the archiver to create redirect responses.

The syntax is quite simple:

::

    # This is a comment.

    # NOTE:  Absolute paths are allowed only at the top of the sources folder.
    /some-path     https://example.com/     301

    # NOTE:  Relative paths are always, and are reinterpreted as relative to the containing folder.
    ./some-path    https://example.com/     302

    # NOTE:  Redirects only for a specific domain.  (The protocol is irelevant.)
    #        (Allowed only at the top of the sources folder.)
    ://example.com/some-path         https://example.com/    303
    http://example.com/some-path     https://example.com/    307
    https://example.com/some-path    https://example.com/    308




Symlinks, hardlinks, loops, and duplicated files
................................................

You freely use symlinks (including pointing outside of the content root) and they will be crawled during archival respecting the "logical" hierarchy they introduce.
(Any loop that you introduce into the hierarchy will be ignored and a warning will be issued.)

You can safely symlink or hardlink the same file (or folder) in multiple places (within the content hierarchy), and its data will be stored only once.
(The same applies to duplicated files that have exactly the same data.)

