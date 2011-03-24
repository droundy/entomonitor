# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=entomonitor

GOFILES=\
	entomonitor.go\

include $(GOROOT)/src/Make.cmd
