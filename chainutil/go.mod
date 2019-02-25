module github.com/endurio/ndrd/chainutil

replace (
	github.com/endurio/ndrd => ../
	github.com/endurio/ndrd/chainutil => ./
)

require (
	github.com/aead/siphash v1.0.1
	github.com/davecgh/go-spew v0.0.0-20171005155431-ecdeabc65495
	github.com/endurio/ndrd v0.0.0-20190213025234-306aecffea32
	github.com/kkdai/bstream v0.0.0-20161212061736-f391b8402d23
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
)
