GOTEST=go test -v
PKGNAME=github.com/lni/goutils

test:
	$(GOTEST) $(PKGNAME)/syncutil
	$(GOTEST) $(PKGNAME)/netutil
