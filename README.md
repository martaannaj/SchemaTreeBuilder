[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) 

## Installation

1. Install the go runtime (and VS Code + Golang tools)
1. Run `go get .` in this folder to install all dependencies
1. Run `go build .` in this folder to build the executable
1. Run `go install .` to install the executable in the $PATH

## Example 

```bash

# This example will assume that you are in the top directory.

# Download a dataset, for example the latest 32GB dataset from wikidata
# curl https://dumps.wikimedia.org/wikidatawiki/entities/latest-truthy.nt.gz --output latest-truthy.nt.gz
# (this example will assume that a dataset called `./testdata/handcrafted.nt` exists)

# Split the dataset for wikidata items and properties
# (TODO: The handcrafted dataset has to be improved  with a better combination of entries)
./SchemaTreeBuilder split-dataset by-prefix ./testdata/handcrafted.nt

# Prepare the dataset and build the Schema Tree (typed variant) (the sort is only required for future 1-in-n splits)
./SchemaTreeBuilder filter-dataset for-schematree ./testdata/handcrafted-item.nt.gz 
gzip -cd ./testdata/handcrafted-item-filtered.nt.gz | sort | gzip > ./testdata/handcrafted-item-filtered-sorted.nt.gz
./SchemaTreeBuilder build-tree-typed ./testdata/handcrafted-item-filtered-sorted.nt.gz


