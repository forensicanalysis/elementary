#!/usr/bin/env python
#  Copyright (c) 2020 Siemens AG
#
#  Permission is hereby granted, free of charge, to any person obtaining a copy of
#  this software and associated documentation files (the "Software"), to deal in
#  the Software without restriction, including without limitation the rights to
#  use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
#  the Software, and to permit persons to whom the Software is furnished to do so,
#  subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included in all
#  copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
#  FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
#  COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
#  IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
#  CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
#  Author(s): Jonas Plum, Korbinian Karl
import json
import logging
import sys

import forensicstore
import storeutil

LOGGER = logging.getLogger(__name__)


class USBForensicStoreExtractor:
    # pylint: disable=fixme,too-few-public-methods

    def __init__(self, store):
        self.artifact_names = ["WindowsUSBDeviceInformations", "WindowsUSBUserMountedDevices",
                               "WindowsUSBVolumeAndDriveMapping"]
        self.forensicstore = store
        # self.filter = store_filter

    def get_usb_usage_data(self):
        """
        Collects all kinds of data about usb devices and their usage on a given forensicstore.

        :return: A list containing dicts about used usb devices
        """
        system_mounted_usb_info = self._get_system_usb_data()
        usb_user_mounted_devices = self._get_user_usb_data(system_mounted_usb_info)
        return self._get_all_usb_data(usb_user_mounted_devices)

    def _get_system_usb_data(self):
        system_mounted_usb_info = {}
        usb_system_mounted_devices = list(self.forensicstore.select([{"artifact": "WindowsUSBVolumeAndDriveMapping"}]))
        if not usb_system_mounted_devices:
            return {}
        usb_system_mounted_devices = [device.get('values') for device in usb_system_mounted_devices].pop()

        # Collects some information about the system mounted usb devices
        for mapping in usb_system_mounted_devices:
            splitted_values = mapping.get("name").split("\\").pop()
            if splitted_values.startswith("Volume"):
                volume_guid = splitted_values.replace("Volume", '')
                try:
                    device_details = bytes.fromhex(mapping.get("data")).decode("utf16")
                    if device_details and device_details.startswith("_??_USBSTOR#Disk"):
                        cleaned_details = device_details.split("&")
                        vendor_name = cleaned_details[1]
                        product_name = cleaned_details[2]
                        usb_revision = cleaned_details[3].split("#")[0]
                        usb_serial_number = cleaned_details[3].split("#")[1]
                        system_mounted_usb_info.update({volume_guid: {"vendor_name": vendor_name.replace("Ven_", ""),
                                                                      "product_name": product_name.replace("Prod_", ""),
                                                                      "usb_revision": usb_revision.replace("Rev_", ""),
                                                                      "usb_uid": usb_serial_number,
                                                                      }})
                except UnicodeError:
                    pass
        return system_mounted_usb_info

    def _get_user_usb_data(self, system_mounted_usb_info: dict):
        usb_user_mounted_data = self.forensicstore.select([{"artifact": "WindowsUSBUserMountedDevices"}])

        # Categorises found system usb usage to user usb usage.
        usb_user_mounted_devices = []
        for device in usb_user_mounted_data:
            splitted_reg_key = device.get("key").split("\\")
            rear_reg_key = splitted_reg_key[-1]
            if rear_reg_key.startswith("{") and rear_reg_key in system_mounted_usb_info:
                usb_dict = system_mounted_usb_info[rear_reg_key]
                usb_dict.update({"volume_guid": rear_reg_key.replace("{", "").replace("}", "")})
                user_sid = splitted_reg_key[1]
                usb_dict.update({"user_sid": user_sid})
                usb_user_mounted_devices.append(usb_dict)

        return usb_user_mounted_devices

    def _get_all_usb_data(self, usb_user_mounted_devices: list):
        usb_device_information = self.forensicstore.select([{"artifact": "WindowsUSBDeviceInformations"}])

        # Combines the gathered information to create a dictionary of actual used usb devices and some meta data.
        # Also keeps track about non mounted usb devices.
        mounted_device_ids = [device.get("usb_uid") for device in usb_user_mounted_devices]
        for device in usb_device_information:
            splitted_reg_key = device.get("key").split("/")

            # TODO: Can probably done better through a regex.
            if len(splitted_reg_key) == 7:
                device_id = splitted_reg_key[-1].split("&")[0]

                if device_id in mounted_device_ids:
                    for user_mounted_device in usb_user_mounted_devices:
                        if device_id == user_mounted_device.get("usb_uid"):
                            friendly_name = [value.get("data") for value in device.get("values")
                                             if value.get("name") == "FriendlyName"].pop()
                            user_mounted_device.update({"friendly_name": friendly_name})
                            user_mounted_device.update(self._get_first_insert_timestamps(device_id))
                else:
                    usb_revision = splitted_reg_key[5].split("&")[3].replace("Rev_", "")
                    device_dict = {"vendor_name": None, "product_name": None, "usb_revision": usb_revision,
                                   "usb_uid": device_id, "volume_guid": None, "user_sid": None}
                    device_dict.update(self._get_first_insert_timestamps(device_id))
                    for value in device.get("values"):
                        if value.get("name") == "FriendlyName":
                            device_dict.update({"friendly_name": value.get("data")})
                    usb_user_mounted_devices.append(device_dict)

        return usb_user_mounted_devices

    def _get_first_insert_timestamps(self, device_id):
        items = self.forensicstore.select([{"name": "setupapi.dev.log"}])
        # fsf = self.forensicstore.remote_fs.open("WindowsDeviceSetup/setupapi.dev.log", mode='rb')
        for item in items:
            if "export_path" not in item:
                continue
            fsf = self.forensicstore.fs.open(item["export_path"], mode='rt')

            inital_timestamp = {"first_insert": None}
            if device_id:

                # Checks each line for the first insert timestamp in setupapi.dev.log
                for line in fsf:
                    try:
                        log_line = line.decode("utf-8")
                        if "Device Install (Hardware initiated)" in log_line and device_id in log_line:
                            splitted_log_line = next(fsf).decode("UTF-8").split(' ')
                            date = splitted_log_line.pop().replace("\r\n", "")
                            time = splitted_log_line.pop()
                            inital_timestamp = {"first_insert": date + " " + time}
                            break

                    except UnicodeDecodeError:
                        pass
        return inital_timestamp

    @staticmethod
    def _get_all_insert_timestamps():
        # TODO get all insert timestamps via the last change timestamp of the usb UID.
        # fs = self.forensicstore
        return []


def main(url):
    LOGGER.debug("process usb")
    store = forensicstore.open(url)

    usb_usage_data = USBForensicStoreExtractor(store).get_usb_usage_data()
    for result in usb_usage_data:
        result["type"] = "usb-device"
        print(json.dumps(result))
    store.close()


if __name__ == '__main__':
    logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
    parser = storeutil.ScriptArgumentParser(
        'usb', description='Process usb artifacts', store_arg=True, filter_arg=False,
    )
    args, _ = parser.parse_known_args(sys.argv[1:])
    main(args.forensicstore)
