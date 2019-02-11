module github.com/endurio/ndrd

replace (
	github.com/endurio/ndrd => ./
	github.com/endurio/ndrd/util => ./util
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/goleveldb v0.0.0-20160330041536-7834afc9e8cd
	github.com/btcsuite/snappy-go v0.0.0-20151229074030-0bdef8d06723 // indirect
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/btcsuite/winsvc v1.0.0
	github.com/davecgh/go-spew v0.0.0-20171005155431-ecdeabc65495
	github.com/endurio/ndrd/util v0.0.0-20180706230648-ab6388e0c60a
	github.com/jessevdk/go-flags v0.0.0-20141203071132-1679536dcc89
	github.com/jrick/logrotate v1.0.0
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
)
