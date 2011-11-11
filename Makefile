GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

TARG=github.com/iNamik/parser.go

GOFILES=\
	impl.go\
	parser.go\
	private.go\

include $(GOROOT)/src/Make.pkg

