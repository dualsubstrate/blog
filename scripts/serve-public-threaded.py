#!/usr/bin/env python3
"""Serve ./public over HTTP for local link checks.

The stdlib `python -m http.server` handler is single-threaded; muffet opens many
parallel connections and requests queue until they hit the crawler timeout. This
server uses ThreadingHTTPServer (Python 3.7+) so concurrent asset requests do not
starve each other.
"""
from __future__ import annotations

import os
import sys
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


def main() -> None:
    if len(sys.argv) != 2:
        print("usage: serve-public-threaded.py <port>", file=sys.stderr)
        sys.exit(2)
    port = int(sys.argv[1], 10)
    root = os.path.dirname(os.path.abspath(__file__))
    public = os.path.join(root, "..", "public")
    os.chdir(public)
    host = "127.0.0.1"
    httpd = ThreadingHTTPServer((host, port), SimpleHTTPRequestHandler)
    httpd.serve_forever()


if __name__ == "__main__":
    main()
