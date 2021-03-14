<h1 align="center">elementary</h1>

<p  align="center">
 <a href="https://godocs.io/github.com/forensicanalysis/elementary"><img src="https://godocs.io/github.com/forensicanalysis/elementary?status.svg" alt="doc" /></a>
</p>

The elementary tool can process forensicstores created with the [artifactcollector](https://github.com/forensicanalysis/artifactcollector).

## 💾 Installation

<details><summary><b>homebrew (macOS and Linux)</b></summary>

If you have the [Homebrew](https://brew.sh/) package manager installed, you can install elementary using:

```bash
brew tap forensicanalysis/tap
brew install elementary
```

</details>
<details><summary><b>scoop (Windows)</b></summary>

If you have the [Scoop](https://scoop.sh/) package manager installed, you can install elementary using:

```bash
scoop bucket add elementary https://github.com/forensicanalysis/homebrew-tap
scoop install elementary
```

</details>
<details><summary><b>deb/rpm (Linux)</b></summary>

Download the .deb or .rpm from the [releases](https://github.com/forensicanalysis/elementary/releases) 
page and install with `dpkg -i` and `rpm -i` respectively.

</details>
<details><summary><b>manually</b></summary>

The GitHub [releases](https://github.com/forensicanalysis/elementary/releases) pages provides binaries for all common systems.

</details>


## 🧑‍💻 Usage

For all commands see `elementary --help`. For all features and flags append `--help` to any command.

<details><summary><b>Unpack a forensicstore</b></summary>

```bash
elementary archive unpack pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

</details>
<details><summary><b>Get connected usb devices</b></summary>

```bash
elementary run usb pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```


</details>
<details><summary><b>Get some autostarts</b></summary>

```bash
elementary run run-keys pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

</details>
<details><summary><b>List installed services</b></summary>

```bash
elementary run services pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

</details>
<details><summary><b>List uninstall entries</b></summary>

```bash
elementary run software pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

</details>
<details><summary><b>List network devices</b></summary>

```bash
elementary run networking pc2dd9f0f_2020-05-16T16-46-25.forensicstore
```

</details>

## 🚫 Limitations

- Most commands only process Windows artifacts
- Prefetch file processing is very slow

## 💬 Contact

For feedback, questions and discussions you can use the [Discussions](https://github.com/forensicanalysis/elementary/discussions) or the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
