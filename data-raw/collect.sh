#!/usr/bin/env bash
find `pwd`/data-raw/top -type f -name "solution_data.json" -exec cat {} \; > `pwd`/data-raw/top.json
