GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

TARG=calc

GOFILES=\
	calc.go\

include $(GOROOT)/src/Make.cmd

