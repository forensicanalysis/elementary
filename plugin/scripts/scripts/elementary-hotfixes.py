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
This plugin parses different registry entries for installed Hotfixes (patches) to the Windows system
as well as to other software components
"""
import json
import logging
import re
import struct
import sys
from collections import defaultdict
from datetime import datetime

import forensicstore
import storeutil

LOGGER = logging.getLogger(__name__)

HOTFIX_PATHS_INSTALLER = [
    'hkey_local_machine\\software\\microsoft\\windows\\currentversion\\component based servicing\\packages\\',
]
HOTFIX_PATHS_ADDITIONAL = [
    'hkey_local_machine\\software\\wow6432node\\microsoft\\updates\\',
    'hkey_local_machine\\software\\microsoft\\updates\\',
]

KB_REGEX = re.compile(r'KB\d+')


def _analyze_installer(obj):
    entries = []
    installer_entries = defaultdict(set)

    hotfix_infos = {v["name"].lower(): v["data"] for v in obj["values"]}

    if hotfix_infos.get('InstallClient') != 'WindowsUpdateAgent':
        return []
    hotfix = KB_REGEX.search(obj["key"].split('\\')[-1])
    if not hotfix:
        # some entries do not have the KB number in the title, but something like "RollupFix", check
        # the InstallLocation value in this case
        location = hotfix_infos.get('InstallLocation')
        if location:
            hotfix = KB_REGEX.search(location)
    if not hotfix:
        LOGGER.info("Non KB entry for WindowsUpdateAgent found: %s", obj["key"])
        return []
    install_high = hotfix_infos.get('InstallTimeHigh')
    install_low = hotfix_infos.get('InstallTimeLow')
    if install_high and install_low:
        timestamp = filetime_to_timestamp(filetime_join(install_high, install_low))
    else:
        timestamp = ''
    installer_entries[hotfix.group(0)].add(timestamp)

    for hotfix in installer_entries:
        entries.append({
            'Hotfix': hotfix,
            'Installed': sorted(installer_entries[hotfix])[0] if installer_entries[hotfix] else '-',
            'Source': 'Component Based Servicing',
            "type": "hotfix"
        })

    return entries


def _analyze_additional(key):
    hotfix = key["key"].split('\\')[-1]
    product = key["key"].split('\\')[-2]
    return [{
        'Hotfix': hotfix,
        'Installed': key["modified_time"],
        'Source': 'Microsoft Updates',
        'Component': product,
        "type": "hotfix"
    }]


def transform(obj):
    if any(map(lambda path: obj["key"].lower().startswith(path), HOTFIX_PATHS_INSTALLER)):
        return _analyze_installer(obj)
    if any(map(lambda path: obj["key"].lower().startswith(path), HOTFIX_PATHS_ADDITIONAL)):
        return _analyze_additional(obj)
    return []


def filetime_join(upper, lower):
    """
    :param upper: upper part of the number
    :param lower: lower part of the number
    """
    return struct.unpack('Q', struct.pack('ii', lower, upper))[0]


def filetime_to_timestamp(filetime_64):
    """
    The FILETIME timestamp is a 64-bit integer that contains the number
    of 100th nano seconds since 1601-01-01 00:00:00.
    The number is usually saved in the registry using two DWORD["values"]
    :return: string of UTC time
    """
    # pylint: disable=invalid-name
    HUNDREDS_OF_NANOSECONDS_IN_A_SECOND = 10000000
    UNIXEPOCH_AS_FILETIME = 116444736000000000

    datetime_stamp = datetime.utcfromtimestamp(
        (filetime_64 - UNIXEPOCH_AS_FILETIME) / HUNDREDS_OF_NANOSECONDS_IN_A_SECOND)
    return datetime_stamp.isoformat()


def main(url):
    LOGGER.debug("search for hotfixes")
    store = forensicstore.open(url)

    hklmsw = "HKEY_LOCAL_MACHINE\\SOFTWARE\\"
    conditions = [
        {'key': hklmsw + "Microsoft\\Windows\\CurrentVersion\\Component Based Servicing\\Packages\\%"},
        {'key': hklmsw + "WOW6432Node\\Microsoft\\Updates\\%\\%"},
        {'key': hklmsw + "Microsoft\\Updates\\%\\%"}
    ]
    for item in store.select(conditions):
        for result in transform(item):
            print(json.dumps(result))
    store.close()


if __name__ == '__main__':
    logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
    parser = storeutil.ScriptArgumentParser(
        'hotfixes', description='Process windows hotfixes', store_arg=True, filter_arg=False,
    )
    args, _ = parser.parse_known_args(sys.argv[1:])
    main(args.forensicstore)
