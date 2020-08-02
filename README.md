<h1 align="center">elementary</h1>

<p  align="center">
 <a href="https://github.com/forensicanalysis/forensicworkflows/actions"><img src="https://github.com/forensicanalysis/elementary/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/elementary"><img src="https://codecov.io/gh/forensicanalysis/elementary/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/elementary"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/elementary" alt="report" /></a>
 <a href="https://pkg.go.dev/github.com/forensicanalysis/elementary"><img src="https://img.shields.io/badge/go.dev-documentation-007d9c?logo=go&logoColor=white" alt="doc" /></a>
 <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Felementary?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Felementary.svg?type=shield"/></a>
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
- Script commands require Python 3.9.0a on Windows

## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
