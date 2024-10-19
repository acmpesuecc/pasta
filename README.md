# pasta

minimal pastebin: written in go 🤌

## usage
```
$ export PASSPHRASE="your-key"
$ export COOLDOWN="time-in-minutes"
$ export LIMIT="number-of-requests-per-cooldown"
$ curl -F "file=@file.txt" "localhost:8080" -H "X-Auth-Passphrase: your-key"
```


### install

```
go install codeberg.org/polarhive/pasta@latest
```

### limitations

I set it up on my server for my personal use. [Ping](https://polarhive.net/contact) me if you wish to use my resources (no spam lol)

> I may delete your files at any time.

---
This repo is hosted on [Codeberg](https://codeberg.org/polarhive/pasta) & mirrored to [GitHub](https://github.com/polarhive/pasta) for traffic.

[![GPL enforced badge](https://img.shields.io/badge/GPL-enforced-blue.svg "This project enforces the GPL.")](https://gplenforced.org)
