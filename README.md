# Whisper
---

Runs an arbitrary command using secrets provided in a [SOPS](https://github.com/mozilla/sops)-encrypted file.

```
$ whisper -f <sops yaml file> [-p] <command>
```

Options:

* `-f` (required) path to a SOPS YAML file
* `-p` (optional) use the rest of the current environment
* `command` (required) the command to provide the decrypted secrets to

SOPS file format just requires a top-level `environment` key:

```
environment:
    MY_PASSWORD: wowthisissecret
    MY_OTHER_KEY: license1234
```

**NOTE:** Absolutely none of this has been vetted so you probably shouldn't actually use it, I was just seeing if I could get it to work. This is literally the first thing I've written in Go.
