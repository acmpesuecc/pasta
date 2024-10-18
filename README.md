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

### docker usage

The docker container can be started either using the `docker run` command or using the `docker-compose.yml` file included.

#### docker run
```
docker build -t pasta:latest .
docker run -d -p 8080:8080 pasta
```

#### docker-compose
```
docker-compose build
docker-compose up
```

### limitations

I set it up on my server for my personal use. [Ping](https://polarhive.net/contact) me if you wish to use my resources (no spam lol)

> I may delete your files at any time.

---
This repo is hosted on [Codeberg](https://codeberg.org/polarhive/pasta) & mirrored to [GitHub](https://github.com/polarhive/pasta) for traffic.

[![GPL enforced badge](https://img.shields.io/badge/GPL-enforced-blue.svg "This project enforces the GPL.")](https://gplenforced.org)
