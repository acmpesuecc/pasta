# pasta

minimal pastebin: written in go 🤌

## usage
```
$ curl -F "file=@file.txt" "localhost:8080"
```

### install

```
go install codeberg.org/polarhive/pasta@latest
```
### systemd service
```sh
# cp pasta.service /etc/systemd/system/pasta.service
$ systemctl enable --now pasta.service
```
If pasta is in the home directory, the pasta-user.service can be used instead, while replacing $USER with the user's username.

### limitations

I set it up on my server for my personal use. [Ping](https://polarhive.net/contact) me if you wish to use my resources (no spam lol)

> I may delete your files at any time.

---
This repo is hosted on [Codeberg](https://codeberg.org/polarhive/pasta) & mirrored to [GitHub](https://github.com/polarhive/pasta) for traffic.

[![GPL enforced badge](https://img.shields.io/badge/GPL-enforced-blue.svg "This project enforces the GPL.")](https://gplenforced.org)
