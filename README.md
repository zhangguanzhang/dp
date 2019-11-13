## demo

<a href="https://asciinema.org/a/QYGxvq2ef43pnhO5HZo2KcHBk?autoplay=1&speed=2"><img src="https://asciinema.org/a/QYGxvq2ef43pnhO5HZo2KcHBk.png" width="836"/></a>

## Usage
Pull the docker images on a machine without docker, and support pulling images from multiple registry at same time
```
$ dp pull
pull all images and write to a tar.gz file without docker daemon.

Usage:
  dp pull [flags]

Aliases:
  pull, p

Examples:

# pull a image or set the name to save
dp pull nginx:alpine
dp pull -o nginx.tar.gz nginx:alpine

# pull image use sha256
dp pull mcr.microsoft.com/windows/nanoserver@sha256:ae443bd9609b9ef06d21d6caab59505cb78f24a725cc24716d4427e36aedabf2

# pull images and set the name to save
dp pull -o project.tar.gz nginx:alpine nginx:1.17.5-alpine-perl

# pull from different registry 
dp pull -o project.tar.gz nginx:alpine gcr.azk8s.cn/google_containers/pause-amd64:3.1


Flags:
  -h, --help              help for pull
  -o, --out-file string   the name will write to
```

## todo

- support quay.io and harbor
- could retry while failed
- multi process download
- with a nice download progress bar
