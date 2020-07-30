#!/usr/bin/env bats
# Copyright (c) 2020 Siemens AG
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# Author(s): Jonas Plum


setup() {
  TESTDIR=$BATS_TMPDIR/bats/$BATS_TEST_NUMBER
  mkdir -p $TESTDIR
}

teardown() {
  rm -rf $TESTDIR
}

@test "run import-json (go)" {
  forensicstore create $TESTDIR/test.forensicstore
  run elementary run import-json --file test/data/import.json $TESTDIR/test.forensicstore --debug
  echo $output
  [ "$status" -eq 0 ]
}

@test "run prefetch (go)" {
  cp -r test/data/example1.forensicstore $TESTDIR/example1.forensicstore
  [ -f "$TESTDIR/example1.forensicstore" ]
  elementary run prefetch $TESTDIR/example1.forensicstore --debug
}

# @test "run plaso (docker)" {
#   cp -r test/data/example1.forensicstore $TESTDIR/example1.forensicstore
#   [ -f "$TESTDIR/example1.forensicstore" ]
#   elementary run plaso $TESTDIR/example1.forensicstore --debug
# }

@test "run export-json (go)" {
  cp -r test/data/example1.forensicstore $TESTDIR/example3.forensicstore
  [ -f "$TESTDIR/example3.forensicstore" ]
  run elementary run export --format jsonl --filter type=file --output $TESTDIR/export.json $TESTDIR/example3.forensicstore --debug
  echo $output
  [ "$status" -eq 0 ]
  [ -f "$TESTDIR/export.json" ]
}
