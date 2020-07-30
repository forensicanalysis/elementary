// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

export const templates = {
  "": {
    icon: 'file-multiple',
    headers: [{
      text: 'Title',
      align: 'left',
      value: 'title',
    },{
      text: 'Type',
      align: 'left',
      value: 'type',
    }]
  },
  file: {
    icon: 'file',
    headers: [{
      text: 'Name',
      align: 'left',
      value: 'name',
    }, {
      text: 'Size',
      value: 'size',
      align: 'right',
    }, {
      text: 'Accessed',
      value: 'atime',
      align: 'right',
    }, {
      text: 'Modified',
      value: 'mtime',
      align: 'right',
    }, {
      text: 'Created',
      value: 'ctime',
      align: 'right',
    }, {
      text: 'Origin Path',
      value: 'origin',
    }],
  },
  'windows-registry-key': {
    icon: 'file-cabinet',
    headers: [{
      text: 'Key',
      align: 'left',
      value: 'key',
    }, {
      text: 'Modified',
      value: 'modified_time',
      align: 'right',
    }],
  },
  process: {
    icon: 'console',
    headers: [{
      text: 'Name',
      align: 'left',
      value: 'name',
    }, {
      text: 'Command Line',
      align: 'left',
      value: 'command_line',
    }, {
      text: 'Created',
      align: 'left',
      value: 'created_time',
    }, {
      text: 'Return Code',
      align: 'left',
      value: 'return_code',
    }],
  },
  info: { icon: 'info' },
  report: { icon: 'file-alt' },
  directory: {
    icon: 'folder',
    headers: [
      {
        text: 'Path',
        value: 'path',
        align: 'left',
      }, {
        text: 'Accessed',
        value: 'atime',
        align: 'right',
      }, {
        text: 'Modified',
        value: 'mtime',
        align: 'right',
      }, {
        text: 'Created',
        value: 'ctime',
        align: 'right',
      },
    ],
  },
  'hotfix': {
    icon: 'bandage',
    headers: [
      {
        text: 'Hotfix',
        value: 'Hotfix',
        align: 'left',
      }, {
        text: 'Source',
        value: 'Source',
        align: 'left',
      }, {
        text: 'Component',
        value: 'Component',
        align: 'left',
      }, {
        text: 'Installed',
        value: 'Installed',
        align: 'right',
      },
    ],
  },
  event: {
    icon: 'chart-timeline',
    headers: [
      {
        text: 'Timestamp',
        value: 'timestamp',
        align: 'left',
      }, {
        text: 'Description',
        value: 'timestamp_desc',
        align: 'left',
      }, {
        text: 'Message',
        value: 'message',
        align: 'left',
      }, {
        text: 'Parser',
        value: 'parser',
        align: 'left',
      },
    ],
  },
  eventlog: {
    icon: 'calendar-clock',
    headers: [
      {
        text: 'Computer',
        value: 'System.Computer',
        align: 'left',
      }, {
        text: 'SystemTime',
        value: 'System.TimeCreated.SystemTime',
        align: 'right',
      }, {
        text: 'EventRecordID',
        value: 'System.EventRecordID',
        align: 'right',
      }, {
        text: 'EventID',
        value: 'System.EventID.Value',
        align: 'right',
      }, {
        text: 'Level',
        value: 'System.Level',
        align: 'right',
      }, {
        text: 'Channel',
        value: 'System.Channel',
        align: 'left',
      }, {
        text: 'Provider',
        value: 'System.Provider.Name',
        align: 'left',
      },
    ],
  },
  'known_network': {
    icon: 'lan',
    headers: [
      {
        text: 'Friendly Name',
        value: 'Friendly Name',
        align: 'left',
      }, {
        text: 'GUID',
        value: 'GUID',
        align: 'left',
      }, {
        text: 'IPs',
        value: 'IPs',
        align: 'right',
      }, {
        text: 'SubNetMask',
        value: 'SubNetMask',
        align: 'right',
      }, {
        text: 'NameServer',
        value: 'NameServer',
        align: 'left',
      }, {
        text: 'DHCP',
        value: 'DHCP',
        align: 'right',
      }, {
        text: 'IP Key Changed',
        value: 'IP Key Changed',
        align: 'right',
      }, {
        text: 'Network Key Changed',
        value: 'Network Key Changed',
        align: 'right',
      },
    ],
  },
  'uninstall_entry': {
    icon: 'delete-sweep',
    headers: [
      {
        text: 'Name',
        value: 'Name',
        align: 'left',
      }, {
        text: 'Version',
        value: 'Version',
        align: 'right',
      }, {
        text: 'Publisher',
        value: 'Publisher',
        align: 'left',
      }, {
        text: 'InstallDate',
        value: 'InstallDate',
        align: 'right',
      }, {
        text: 'Source',
        value: 'Source',
        align: 'left',
      }, {
        text: 'Location',
        value: 'Location',
        align: 'left',
      },
    ],
  },
  'prefetch': {
    icon: 'reload',
    headers: [
      {
        text: 'Executable',
        value: 'Executable',
        align: 'left',
      }, {
        text: 'Last Run Times',
        value: 'LastRunTimes',
        align: 'left',
      }, {
        text: 'RunCount',
        value: 'RunCount',
        align: 'right',
      }
    ],
  },
};
