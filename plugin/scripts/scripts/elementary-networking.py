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
import json
import sys

import forensicstore
import storeutil

NAMES_KEY = r"\control\network\{4d36e972-e325-11ce-bfc1-08002be10318}"
INTERFACE_KEY = r'\services\tcpip\parameters\interfaces'


def transform(objs):
    """
    This plugin parses network configuration from the SYSTEM hive and correlates it with other keys to
    get the network interface name and the last configuration
    """
    names_keys = []
    interface_keys = []
    for obj in objs:
        key = obj["key"].lower()
        if NAMES_KEY in key:
            names_keys.append(obj)
        elif INTERFACE_KEY in key:
            interface_keys.append(obj)

    results = []
    for interface_key in interface_keys:
        result = {}

        interface_id = interface_key["key"].split('\\')[-1]

        interface_infos = {}

        # create a map from values
        if "values" in interface_key:
            # return interface_key["values"]
            interface_infos = {
                v["name"].lower(): v["data"] if "data" in v else ""
                for v in interface_key["values"]
            }

        missing_value = '(not specified)'
        result['type'] = "known_network"
        result['GUID'] = interface_id
        result['DHCP'] = 'yes' if interface_infos.get(
            'enabledhcp', missing_value) == 1 else 'no'
        result['IPs'] = interface_infos.get('ipaddress', missing_value)
        result['SubNetMask'] = interface_infos.get('subnetmask',
                                                   missing_value)
        result['NameServer'] = interface_infos.get('nameserver',
                                                   missing_value)
        result['IP Key Changed'] = interface_key['modified_time']

        # search names key
        name_key = None
        for key in names_keys:
            if interface_id.lower() in key["key"].lower():
                name_key = key

        if name_key is not None:
            # create a map from values
            name_infos = {
                v["name"].lower(): v["data"] if "data" in v else ""
                for v in name_key["values"]
            }

            result['Network Key Changed'] = name_key["modified_time"]
            result['Friendly Name'] = name_infos["name"]
        else:
            result['Network Key Changed'] = missing_value
            result['Friendly Name'] = missing_value
        results.append(result)

    return results


def main(url):
    store = forensicstore.open(url)
    conditions = [
        {'key': r"HKEY_LOCAL_MACHINE\SYSTEM\%ControlSet%\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}\%"},
        {'key': r"HKEY_LOCAL_MACHINE\SYSTEM\%ControlSet%\Services\Tcpip\Parameters\Interfaces\%"}
    ]
    items = store.select(conditions)
    for result in transform(items):
        print(json.dumps(result))
    store.close()


if __name__ == '__main__':
    parser = storeutil.ScriptArgumentParser(
        'networking', description='Process networking artifacts', store_arg=True, filter_arg=False,
    )
    args, _ = parser.parse_known_args(sys.argv[1:])
    main(args.forensicstore)
