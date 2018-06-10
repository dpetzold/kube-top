kube-top
========

[![Build Status](https://circleci.com/gh/dpetzold/kube-top.svg?style=svg)](https://circleci.com/gh/dpetzold/kube-top)

Aggregates and provides visual representation of the following kubectl commands:

* `kubectl top pods`
* `kubectl top nodes`
* `kubectl get pods`
* `kubectl get events`

<img src="./_img/demo.gif" width="80%">

## Usage

* `-namespace` - Override the default namespace. For all namespaces use "".

## Run
```
docker run -it \
    -v/home/derrick/.kube:/.kube \
    -v/home/derrick/.config/gcloud:/.config/gcloud \
    -v/etc/ssl/certs:/etc/ssl/certs \
    dpetzold/kube-top
```

## Build
```
go get github.com/dpetzold/kube-top/cmd/kube-top
```
