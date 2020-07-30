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

mkdir -p test/data

rm -rf test/data/example1.forensicstore
if [ ! -f "example1.forensicstore" ]; then
    curl --fail --output example1.forensicstore --location https://download.artifacthub.org/forensics/example1.forensicstore
fi
mv example1.forensicstore test/data

rm -rf test/data/example2.forensicstore
if [ ! -f "example2.forensicstore" ]; then
    curl --fail --output example2.forensicstore --location https://download.artifacthub.org/forensics/example2.forensicstore
fi
mv example2.forensicstore test/data

rm -rf test/data/usb.forensicstore
if [ ! -f "usb.forensicstore" ]; then
    curl --fail --output usb.forensicstore --location https://download.artifacthub.org/forensics/usb.forensicstore
fi
mv usb.forensicstore test/data

if [ ! -f "win10_mock.zip" ]; then
    curl --fail --output win10_mock.zip --location https://download.artifacthub.org/windows/win10_mock.zip
fi
unzip win10_mock.zip
mv win10_mock.vhd test/data
