#!/bin/bash
date
../dev/testdata/seed.sh
# todo usda import
LOG_LEVEL=info ./bin/gourd sync
date
