so9p
====

Sort of a 9p protocol, but much has been cleaned up.
It's reasonably fast. I tested with RSC's fuse code but lost that bit!

go get github.com/rminnich/so9p/so9p
go install github.com/rminnich/so9p/so9p
go get github.com/rminnich/so9p/so9ptest
go install github.com/rminnich/so9p/so9ptest

so9ptest -s=true [-dbg=true] &
so9ptest -s=false [-dbg=true] [file on host to try reading e.g. /etc/hosts]
