CMDS=tradfri tradfrid

.DEFAULT_GOAL := all

all: $(CMDS)

%: cmd/%/main.go
	go build -o $@ $<

.PHONY: clean
clean:
	rm $(CMDS)
