build:
	govendor build +local

install-libs:
	govendor install +vendor,^program
