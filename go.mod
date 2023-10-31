module easylog

go 1.20

require (
	github.com/fsnotify/fsnotify v1.4.3-0.20161026203122-fd9ec7deca8b
	github.com/gookit/color v1.5.4
	github.com/papertrail/go-tail v0.0.0-20221103124010-5087eb6a0a07
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	golang.org/x/crypto v0.14.0
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	gopkg.in/yaml.v2 v2.4.0
	nextworks.it/libosr/go v0.1.0
	nextworks.it/nxw-golibrary v0.1.0
	nxw v0.1.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.13.0 // indirect
)


replace nextworks.it/libosr/go v0.1.0 => ./submodules/libosr/go/
replace nxw v0.1.0 => ./submodules/nxw-grpc/built/go/nxw/
replace nextworks.it/nxw-golibrary v0.1.0 => ./submodules/nxw-golibrary/
//workaround per forzare la versione di grpc alla 1.27.0
replace google.golang.org/grpc => google.golang.org/grpc v1.27.0
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.6
replace go.etcd.io/etcd => go.etcd.io/etcd v3.3.12+incompatible