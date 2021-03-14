#!/usr/bin/env python
# Copyright (c) 2019 Siemens AG
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
# Author(s): Demian Kellermann

"""
This plugin parses the Windows Uninstaller registry keys for a list of installed software
"""
import json
import sys

import forensicstore
import storeutil


def transform(obj):
    uninstall_entry = {
        "Name": "",
        "Key Timestamp": obj.get("modified", ""),
        "Version": "",
        "Publisher": "",
        "InstallDate": "",
        "Source": "",
        "Location": "",
        "Uninstall": "",
        "Key": obj["key"]
    }

    if "values" not in obj:
        return []

    uninstall_infos = {v["name"].lower(): v["data"] if "data" in v else "" for v in obj["values"]}

    uninstall_entry['Name'] = uninstall_infos.get("displayname", '')
    uninstall_entry['Version'] = uninstall_infos.get("displayversion", '')
    uninstall_entry['Publisher'] = uninstall_infos.get("publisher", '')
    if uninstall_infos.get('installdate', None) is not None:
        strnum = str(uninstall_infos.get('installdate'))
        uninstall_entry['InstallDate'] = '{}-{}-{}'.format(
            strnum[0:4], strnum[4:6], strnum[6:8])
    else:
        uninstall_entry['InstallDate'] = ''
    uninstall_entry['Source'] = uninstall_infos.get('installsource', '')
    uninstall_entry['Location'] = uninstall_infos.get(
        'installlocation', '')
    uninstall_entry['Uninstall'] = uninstall_infos.get(
        'uninstallstring', '')
    uninstall_entry["type"] = "uninstall_entry"

    return [uninstall_entry]


def main(url):
    store = forensicstore.open(url)
    conditions = [{
        'key': "HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\%"
    }, {
        'key': "HKEY_USERS\\%\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\%"
    }]
    for item in store.select(conditions):
        results = transform(item)
        for result in results:
            print(json.dumps(result))
    store.close()


if __name__ == '__main__':
    parser = storeutil.ScriptArgumentParser(
        'software', description='Process uninstall entries', store_arg=True, filter_arg=False,
    )
    args, _ = parser.parse_known_args(sys.argv[1:])
    main(args.forensicstore)
