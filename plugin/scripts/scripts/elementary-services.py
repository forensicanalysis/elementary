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

""" Windows services plugin """
import json
import sys

import forensicstore
import storeutil


def transform(objs):
    """
    This plugin parses the Windows SYSTEM hive file for all installed services and displays their startup state
    and other parameters
    """
    # pylint: disable=too-many-locals

    params_for_services = {}
    services = []
    for obj in objs:
        key = obj["key"].lower()
        if key.endswith('\\parameters'):
            params_for_services[key.replace('\\parameters', '')] = obj
        elif key.split('\\')[-2] != 'services':
            # skip other subkeys
            continue
        else:
            services.append(obj)

    parse_result = []
    for obj in services:
        parameters_key = params_for_services.get(obj["key"].lower(), None)
        if parameters_key and "values" in parameters_key:
            parameters = {v["name"].lower(): v["data"] for v in parameters_key["values"] if "data" in v and "name" in v}
            parameters_last_written = parameters_key["modified_time"]
        else:
            parameters = {}
            parameters_last_written = ''

        if "values" not in obj:
            continue

        service_values = {v["name"].lower(): v["data"] for v in obj["values"] if "data" in v and "name" in v}
        start_mode = int(service_values.get('start', 0x4))
        service_type = service_values.get('type', None)

        service = {
            "type": "service",
            'Name': obj["key"].split("\\")[-1],
            'Description': service_values.get('description', ''),
            'DisplayName': service_values.get('displayname', ''),
            'Group': service_values.get('group', ''),
            'Start Mode': STARTUP_TYPES.get(start_mode, start_mode),
            'Service Type': _servicetype_from_bitmask(service_type),
            'ImagePath': service_values.get('imagepath', ''),
            'Service DLL': parameters.get('servicedll'),
            'Service Key Changed': obj["modified_time"],
            'Parameters Key Changed': parameters_last_written,
            'Source Key': obj["key"]
        }

        parse_result.append(service)

    return parse_result


STARTUP_TYPES = {
    # https://docs.microsoft.com/de-de/dotnet/api/system.serviceprocess.servicestartmode?view=netframework-4.7.2
    0x0: "Boot",
    0x1: "System",
    0x2: "Automatic",
    0x3: "Manual",
    0x4: "Disabled"
}

# Enhancement: When a service is configured for a delayed automatic start, a DWORD registry value is present
# under HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\services\<service name>. It's named DelayedAutoStart and
# it's set to 1, along with a Start value of 2 — which is the normal value for automatic start.

SERVICE_TYPES = {
    # https://docs.microsoft.com/de-de/dotnet/api/system.serviceprocess.servicetype?view=netframework-4.7.2
    0x1: "KernelDriver",
    0x2: "FilesystemDriver",
    0x4: "Adapter",
    0x8: "RecognizerDriver",
    0x10: "Win32OwnProcess",
    0x20: "Win32ShareProcess",
    0x100: "InteractiveProcess",
}


def _servicetype_from_bitmask(mask):
    if not mask:
        return ""
    mask = int(mask)
    types = []
    for bit in SERVICE_TYPES:
        if bit & mask:
            types.append(SERVICE_TYPES[bit])
    return ', '.join(types)


def main(url):
    print(json.dumps({
        "header": ['Name', 'Description', 'DisplayName', 'Group', 'Start Mode',
                   'Service Type', 'ImagePath', 'Service DLL',
                   'Service Key Changed', 'Parameters Key Changed', 'Source Key']}))

    store = forensicstore.open(url)
    conditions = [{'key': "HKEY_LOCAL_MACHINE\\SYSTEM\\%ControlSet%\\Services\\%"}]
    items = list(store.select(conditions))
    for result in transform(items):
        print(json.dumps(result))
    store.close()


if __name__ == '__main__':
    parser = storeutil.ScriptArgumentParser(
        'services', description='Process windows services', store_arg=True, filter_arg=False,
    )
    args, _ = parser.parse_known_args(sys.argv[1:])
    main(args.forensicstore)
