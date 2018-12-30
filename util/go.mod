module github.com/endurio/ndrd/util

replace (
	github.com/endurio/ndrd => ../
	github.com/endurio/ndrd/util => ./
)

require (
	github.com/aead/siphash v1.0.1
	github.com/davecgh/go-spew v0.0.0-20171005155431-ecdeabc65495
	github.com/endurio/ndrd v0.0.0-20181229112439-ce9c0a3f5f31
	github.com/kkdai/bstream v0.0.0-20161212061736-f391b8402d23
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
)
