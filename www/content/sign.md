---
title: Signing
series: customization
hideFromIndex: true
weight: 60
---

GoReleaser can sign some or all of the generated artifacts. Signing ensures
that the artifacts have been generated by yourself and your users can verify
that by comparing the generated signature with your public signing key.

Signing works in combination with checksum files and it is generally sufficient
to sign the checksum files only.

The default is configured to create a detached signature for the checksum files
with [GnuPG](https://www.gnupg.org/) and your default key. To enable signing
just add

```yaml
# goreleaser.yml
signs:
  - artifacts: checksum
```

To customize the signing pipeline you can use the following options:

```yml
# .goreleaser.yml
signs:
  -
    # name of the signature file.
    # '${artifact}' is the path to the artifact that should be signed.
    #
    # defaults to `${artifact}.sig`
    signature: "${artifact}_sig"

    # path to the signature command
    #
    # defaults to `gpg`
    cmd: gpg2

    # command line arguments for the command
    #
    # to sign with a specific key use
    # args: ["-u", "<key id, fingerprint, email, ...>", "--output", "${signature}", "--detach-sign", "${artifact}"]
    #
    # defaults to `["--output", "${signature}", "--detach-sign", "${artifact}"]`
    args: ["--output", "${signature}", "${artifact}"]


    # which artifacts to sign
    #
    #   checksum: only checksum file(s)
    #   all:      all artifacts
    #   none:     no signing
    #
    # defaults to `none`
    artifacts: all
```