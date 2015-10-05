Gocov-Merge
===========

Merge coverage reports from gocov test into a single number.

Installation
------------

  go get github.com/dominikschulz/gocov-merge

Usage
-----

  gocov test ./... | gocov-merge

Purpose
-------

This is but a tiny utility to merge the ouput for gocov test into a single number until
gocov issue #68 is eventually implemented.

