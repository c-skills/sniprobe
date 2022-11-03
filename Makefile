GO=go
CXX=c++
CXXFLAGS=-Wall -O2 -pedantic -std=c++11

SSL_LIB  = "/usr/lib/x86_64-linux-gnu/libssl.so.1.1"
GOOD_SNI = "iranpress.com"

DEFS=-DSSL_LIB=\"$(SSL_LIB)\" -DGOOD_SNI=\"$(GOOD_SNI)\"

all: sniprobe sniswitcher.so

clean:
	rm -f sniprobe sniswitcher.so

sniprobe: sniprobe.go
	$(GO) build sniprobe.go

sniswitcher.so: sniswitcher.cc
	$(CXX) -c -fPIC $(CXXFLAGS) $(DEFS) sniswitcher.cc
	$(CXX) -shared -Wl,-soname=sniswitcher sniswitcher.o -ldl -o sniswitcher.so

