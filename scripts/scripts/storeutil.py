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

import argparse


def merge_conditions(list_a, list_b):
    if list_a is None:
        return list_b
    if list_b is None:
        return list_a
    list_c = []
    for item_a in list_a:
        for item_b in list_b:
            list_c.append({**item_a, **item_b})
    return list_c


class DictListAction(argparse.Action):
    # pylint: disable=too-few-public-methods

    def __call__(self, parser, namespace, values, option_string=None):
        flag = {}
        for kv in values.split(","):
            key, value = kv.split("=")
            flag[key] = value
        if hasattr(namespace, self.dest):
            flags = getattr(namespace, self.dest)
            if flags is not None:
                flags.append(flag)
                setattr(namespace, self.dest, flags)
                return
        setattr(namespace, self.dest, [flag])


class ScriptArgumentParser(argparse.ArgumentParser):
    def __init__(self, subcommand, store_arg, filter_arg, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.subcommand = subcommand
        if store_arg:
            self.add_argument('forensicstore', type=str, help='the processed forensicstore')
        if filter_arg:
            self.add_argument(
                "--filter",
                dest="filter",
                action=DictListAction,
                metavar="type=file,name=System.evtx...",
                help="filter processed items")
