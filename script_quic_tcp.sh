#!/bin/bash
export GOOS=linux
export VOLANTMQ_WORK_DIR=/usr/lib/volantmq
export VOLANTMQ_BUILD_FLAGS="-i"
export VOLANTMQ_PLUGINS_DIR=/usr/lib/volantmq/plugins
export GO111MODULE=on
export GOPATH=/home/ffp27/go

sudo mkdir -p $VOLANTMQ_WORK_DIR/bin
sudo mkdir -p $VOLANTMQ_WORK_DIR/conf
sudo mkdir -p $VOLANTMQ_PLUGINS_DIR

#Build the broker
export GO111MODULE=off
go get github.com/fatimafp95/quic-go
cd $GOPATH/src/github.com/fatimafp95/quic-go
export GO111MODULE=on
go mod tidy
export GO111MODULE=off
go get github.com/fatimafp95/volantmq_2
cd $GOPATH/src/github.com/fatimafp95/volantmq_2
sed -i 's@administrator@ffp27@g' go.mod
export GO111MODULE=on
go mod tidy
sudo cp $GOPATH/src/github.com/fatimafp95/volantmq_2/tools/print_version.sh /bin
sudo chmod +x /bin/print_version.sh
export GO111MODULE=off
go get github.com/troian/govvv
cd cmd/volantmq/
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -o volantmq
sudo cp volantmq $VOLANTMQ_WORK_DIR/bin/volantmq

#Debug plugin
go get gitlab.com/VolantMQ/vlplugin/debug
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/debug
git checkout tags/v0.0.8
export GO111MODULE=on
go mod tidy
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/debug.so

#Health plugin
export GO111MODULE=off
go get gitlab.com/VolantMQ/vlplugin/health
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/health
git checkout tags/v0.0.8
export GO111MODULE=on
go mod tidy
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/health.so

#Metrics plugins 
export GO111MODULE=off
go get gitlab.com/VolantMQ/vlplugin/monitoring/prometheus
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/monitoring/prometheus
git checkout tags/v0.0.5
export GO111MODULE=on
go mod tidy
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/prometheus.so

export GO111MODULE=off
go get gitlab.com/VolantMQ/vlplugin/monitoring/systree
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/monitoring/systree
git checkout tags/v0.0.7
export GO111MODULE=on
go mod tidy
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/systree.so

#Persistence plugins
export GO111MODULE=off
go get gitlab.com/VolantMQ/vlplugin/persistence/bbolt
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/persistence/bbolt
git checkout tags/v0.0.7
export GO111MODULE=on
go mod tidy
cd plugin
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/persistence_bbolt.so

#Auth plugins
export GO111MODULE=off
go get gitlab.com/VolantMQ/vlplugin/auth/http
cd $GOPATH/src/gitlab.com/VolantMQ/vlplugin/auth/http
git checkout tags/v0.0.2
GO111MODULE=on
go mod tidy
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -buildmode=plugin -ldflags "-X main.version=$(print_version.sh)" -o $VOLANTMQ_WORK_DIR/plugins/auth_http.so

#Running the broker with config.yaml file
cd $GOPATH/src/github.com/fatimafp95/volantmq_2
sudo cp -r /usr/lib/volantmq /var/lib
#sudo /usr/lib/volantmq/bin/volantmq --config=examples/config.yaml
cd cmd/volantmq
sudo ./volantmq --config=../../examples/config.yaml

#................................................
#MQTT/TCP
#Build and run as MQTT+TCP
export GO111MODULE=off
go get github.com/VolantMQ/volantmq
cd $GOPATH/src/github.com/VolantMQ/volantmq
git checkout tags/v0.4.0-rc.6 
export GO111MODULE=on
go mod tidy
sudo cp $GOPATH/src/github.com/VolantMQ/volantmq/tools/print_version.sh /bin
sudo chmod +x /bin/print_version.sh
export GO111MODULE=off
go get github.com/troian/govvv
cd cmd/volantmq/
sudo /usr/local/go/bin/go build $VOLANTMQ_BUILD_FLAGS -o volantmq
sudo cp volantmq $VOLANTMQ_WORK_DIR/bin/volantmq
cd $GOPATH/src/github.com/VolantMQ/volantmq
sudo cp -r /usr/lib/volantmq /var/lib
cd cmd/volantmq
sudo ./volantmq --config=../../examples/config.yaml

