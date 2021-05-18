#!/usr/bin/env bash

go build .

./SchemaTreeBuilder build-tree-typed --number-pointers=3 ./testdata/handcrafted-item-filtered-sorted.nt.gz