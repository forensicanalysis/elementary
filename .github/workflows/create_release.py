#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import hashlib
import json
import sys
import urllib.request

homebrew_template = """class Elementary < Formula
  desc "🕵️ Process and show forensic artifacts in forensicstores"
  homepage "https://forensicanalysis.github.io/documentation"
  url "%s"
  sha256 "%s"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = HOMEBREW_CACHE/"go_cache"
    (buildpath/"src/github.com/forensicanalysis/elementary").install buildpath.children

    cd "src/github.com/forensicanalysis/elementary" do
      system "go", "build", "-o", bin/"elementary", "main.go"

      # system bin/"elementary", "install", "-f"
      # prefix.install_metafiles
    end
  end

  test do
    store = testpath/"hops-yeast-malt-water.forensicstore"
    system "#\{bin\}/elementary", "create", store
    assert_predicate testpath/store, :exist?
  end
end
"""

scoop_template = """{
  "version": "%s",
  "architecture": {
    "64bit": {
      "url": "%s",
      "bin": [
        "elementary.exe"
      ],
      "hash": "%s"
    }
  },
  "homepage": "https://forensicanalysis.github.io/documentation",
  "license": "MIT",
  "description": "Process and show forensic artifacts (e.g. eventlogs, usb devices, network devices...) in forensicstores"
}
"""

with open(sys.argv[1]) as io, open("Formula/elementary.rb", "w+") as homebrew, open("elementary.json", "w+") as scoop:
    release = json.load(io)

    with urllib.request.urlopen(release["tarball_url"]) as f:
        tarball = f.read()
        sha256 = hashlib.sha256(tarball).hexdigest()

        homebrew.write(homebrew_template % (release["tarball_url"], sha256))

    with urllib.request.urlopen(sys.argv[2]) as f:
        zipfile = f.read()
        sha256 = hashlib.sha256(zipfile).hexdigest()
        scoop.write(scoop_template % (release["tag_name"][1:], sys.argv[2], sha256))
