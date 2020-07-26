## demo

推荐去使用[skeopeo](https://github.com/containers/skopeo)，此项目已经不维护和造轮子，如果有兴趣可以借鉴下思路，azk8s.cn已经失效，如果借鉴过程可以去掉azk8s的代理

<a href="https://asciinema.org/a/QYGxvq2ef43pnhO5HZo2KcHBk?autoplay=1&speed=2"><img src="https://asciinema.org/a/QYGxvq2ef43pnhO5HZo2KcHBk.png" width="836"/></a>

## Usage

### check

Check if the images belongs to scheme2.Manifest
```
$ dp c dduportal/bats:0.4.0
scheme2.Manifest: []
scheme1.Manifest: [dduportal/bats:0.4.0]
$ dp c dduportal/bats:0.4.0 nginx:alpine --only
scheme2.Manifest: [nginx:alpine]
```

### pull

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
## attention

Many of the images on quay.io are still scheme1.manifest, and some of the mirror images of some other domain names are long ago. These images will not be successfully pulled.

## todo

- could retry while failed
- multi process download
- with a nice download progress bar
- support quay.io and harbor(hard for quay.io!!)
