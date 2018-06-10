GIT_COMMIT = $(shell git rev-parse --short HEAD)
GO_SOURCE_FILES = $(shell find pkg -type f -name "*.go")


build: $(GO_SOURCE_FILES)
	go build -i -ldflags "-X main.GitCommit=${GIT_COMMIT} -extldflags '-static'" -o kube-top ./cmd/kube-top


vendor: glide.yaml
	glide up -v


docker-build:
	docker build --rm --no-cache -t dpetzold/kube-top .
	docker push dpetzold/kube-top

run:
	docker run --rm -it \
		-v${HOME}/.kube:/.kube \
		-v${HOME}/.config/gcloud:/.config/gcloud \
		-v/etc/ssl/certs:/etc/ssl/certs \
		dpetzold/kube-top
