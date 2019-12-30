CMDS=tradfri tradfrid
PROTO=api/remote.pb.go

tradfri: $(PROTO)
tradfrid: $(PROTO)

.DEFAULT_GOAL := all

all: $(CMDS)

%.pb.go: %.proto
	protoc --go_out=. $<

%: cmd/%/main.go
	go build -o $@ $<

.PHONY: clean
clean:
	rm -f $(CMDS) $(PROTO)
