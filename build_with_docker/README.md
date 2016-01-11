build golang binary with docker
===

## build files

- build.sh
- Dockerfile
- go_build.sh

### build docker image

- pull golang docker image
- install tools (glide, gox)
- pull golang repository

### go build on docker

- launch docker image that created above
- change branch
- run go_build.sh(install packages with glide and build with gox)
- store binary to bin directory that has mounted on docker host machine.
