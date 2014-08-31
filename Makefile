.PHONY: all rhel7

all: vrouter

GPATH = $(PWD)/build
GOPATH := $(GPATH):${GOPATH}
GOBIN = 
export GOPATH

URL = github.com/zhgwenming
REPO = vrouter

URLPATH = $(GPATH)/src/$(URL)

vrouter: cmd/vrouter/*.go
	@[ -d $(URLPATH) ] || mkdir -p $(URLPATH)
	@ln -nsf $(PWD) $(URLPATH)/$(REPO)
	go install $(URL)/$(REPO)/cmd/vrouter

rhel7: galerabalancer

galerabalancer: *.go
	go build -compiler gccgo -o $@

clean:
	rm -fv build/bin/*
	rm -fv lb cmd/vrouter/vrouter vrouter
