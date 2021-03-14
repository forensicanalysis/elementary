<h1 align="center">elementary</h1>

<p  align="center">
 <a href="https://godocs.io/github.com/forensicanalysis/elementary"><img src="https://godocs.io/github.com/forensicanalysis/elementary?status.svg" alt="doc" /></a>
</p>

The elementary tool can process forensicstores created with the [artifactcollector](https://github.com/forensicanalysis/artifactcollector).

## Installation

Just get the binary:

### [💾 Download](https://github.com/forensicanalysis/elementary/releases)

## Usage

For all commands see `elementary --help`. For all features and flags append `--help` to any command.

### Unpack a forensicstore

```bash
elementary archive unpack pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

### Process a forensicstore

#### Get connected usb devices

```bash
elementary run usb pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

#### Get some autostarts

```bash
elementary run run-keys pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

#### List installed services

```bash
elementary run services pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

#### List uninstall entries

```bash
elementary run software pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

#### List network devices

```bash
elementary run networking pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

## Limitations

- Most commands only process Windows artifacts
- Prefetch file processing is very slow

## Contact

For feedback, questions and discussions you can use the [Discussions](https://github.com/forensicanalysis/elementary/discussions) or the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
