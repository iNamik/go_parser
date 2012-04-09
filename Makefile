GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

TARG=github.com/iNamik/go_parser

GOFILES=\
	impl.go\
	parser.go\
	private.go\

include $(GOROOT)/src/Make.pkg

