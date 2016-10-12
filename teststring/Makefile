all: str

str: runtime
	go install -a -v -x -work -buildmode=shared -linkshared ./
	ls -l ~/go/pkg/linux_amd64_dynlink/libgithub.com-xiaonanln-vacuum-teststring.so
	nm ~/go/pkg/linux_amd64_dynlink/libgithub.com-xiaonanln-vacuum-teststring.so > teststring.nm

runtime:
	$(MAKE) rmruntime
	go install -buildmode=shared runtime sync/atomic unicode/utf8 syscall time reflect strconv fmt io math os log
	$(MAKE) lsruntime

lsruntime:
	ls -l /usr/local/go/pkg/linux_amd64_dynlink/

rmruntime:
	rm -rf /usr/local/go/pkg/linux_amd64_dynlink/*.so

