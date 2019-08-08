.. _licence_rst:

Licences
========

The Monet Toolchain is an open source project, licensed under the `MIT License
<https://opensource.org/licenses/MIT>`__ (`TLDR version
<https://tldrlegal.com/license/mit-license>`__). The software is provided as-is
and we are not liable. We use many other libraries to build the Toolchain. This
section presents the output of `Glice <https://github.com/ribice/glice>`__ for
the Monet Toolchain. Glice reports on the licences used within a golang
project.

The 3 tables are for the `monetd <https://github.com/mosaicnetworks/monetd>`__,
`EVM-Lite <https://github.com/mosaicnetworks/evm-lite>`__ and `Babble
<https://github.com/mosaicnetworks/babble>`__ repositories respectively. These
tables are the output from ``glice -r``, which only looks one level deep. These
tables are for information only and are not legal advice.


``monetd``:

..

    +---------------------------------------------+-----------------------------------------------------+--------------+
    |                 DEPENDENCY                  |                       REPOURL                       |   LICENSE    |
    +---------------------------------------------+-----------------------------------------------------+--------------+
    | github.com/AndreasBriese/bbloom             | https://github.com/AndreasBriese/bbloom             | Other        |
    | github.com/btcsuite/btcd                    | https://github.com/btcsuite/btcd                    | isc          |
    | github.com/dgraph-io/badger                 | https://github.com/dgraph-io/badger                 | Apache-2.0   |
    | github.com/dgryski/go-farm                  | https://github.com/dgryski/go-farm                  | Other        |
    | github.com/docker/docker                    | https://github.com/docker/docker                    | Apache-2.0   |
    | github.com/ethereum/go-ethereum             | https://github.com/ethereum/go-ethereum             | LGPL-3.0     |
    | github.com/fatih/color                      | https://github.com/fatih/color                      | MIT          |
    | github.com/fsnotify/fsnotify                | https://github.com/fsnotify/fsnotify                | bsd-3-clause |
    | github.com/golang/protobuf                  | https://github.com/golang/protobuf                  | bsd-3-clause |
    | github.com/gorilla/mux                      | https://github.com/gorilla/mux                      | bsd-3-clause |
    | github.com/hashicorp/hcl                    | https://github.com/hashicorp/hcl                    | MPL-2.0      |
    | github.com/magiconair/properties            | https://github.com/magiconair/properties            | Other        |
    | github.com/mattn/go-colorable               | https://github.com/mattn/go-colorable               | MIT          |
    | github.com/mattn/go-isatty                  | https://github.com/mattn/go-isatty                  | MIT          |
    | github.com/mgutz/ansi                       | https://github.com/mgutz/ansi                       | MIT          |
    | github.com/mitchellh/mapstructure           | https://github.com/mitchellh/mapstructure           | MIT          |
    | github.com/mosaicnetworks/babble            | https://github.com/mosaicnetworks/babble            | MIT          |
    | github.com/mosaicnetworks/evm-lite          | https://github.com/mosaicnetworks/evm-lite          | MIT          |
    | github.com/pelletier/go-toml                | https://github.com/pelletier/go-toml                | MIT          |
    | github.com/pkg/errors                       | https://github.com/pkg/errors                       | bsd-2-clause |
    | github.com/sirupsen/logrus                  | https://github.com/sirupsen/logrus                  | MIT          |
    | github.com/spf13/afero                      | https://github.com/spf13/afero                      | Apache-2.0   |
    | github.com/spf13/cast                       | https://github.com/spf13/cast                       | MIT          |
    | github.com/spf13/cobra                      | https://github.com/spf13/cobra                      | Apache-2.0   |
    | github.com/spf13/jwalterweatherman          | https://github.com/spf13/jwalterweatherman          | MIT          |
    | github.com/spf13/pflag                      | https://github.com/spf13/pflag                      | bsd-3-clause |
    | github.com/spf13/viper                      | https://github.com/spf13/viper                      | MIT          |
    | github.com/ugorji/go                        | https://github.com/ugorji/go                        | MIT          |
    | github.com/x-cray/logrus-prefixed-formatter | https://github.com/x-cray/logrus-prefixed-formatter | MIT          |
    | golang.org/x/crypto/ssh/terminal            | https://go.googlesource.com/crypto                  |              |
    | golang.org/x/net/internal/timeseries        | https://go.googlesource.com/net                     |              |
    | golang.org/x/net/trace                      | https://go.googlesource.com/net                     |              |
    | golang.org/x/sys/unix                       | https://go.googlesource.com/sys                     |              |
    | golang.org/x/text/transform                 | https://go.googlesource.com/text                    |              |
    | golang.org/x/text/unicode/norm              | https://go.googlesource.com/text                    |              |
    | gopkg.in/yaml.v2                            |                                                     |              |
    +---------------------------------------------+-----------------------------------------------------+--------------+

``EVM-Lite``:

..

    +----------------------------------+-----------------------------------------+----------+
    |            DEPENDENCY            |                 REPOURL                 | LICENSE  |
    +----------------------------------+-----------------------------------------+----------+
    | github.com/sirupsen/logrus       | https://github.com/sirupsen/logrus      | MIT      |
    | golang.org/x/crypto/ssh/terminal | https://go.googlesource.com/crypto      |          |
    | golang.org/x/sys/unix            | https://go.googlesource.com/sys         |          |
    | github.com/ethereum/go-ethereum  | https://github.com/ethereum/go-ethereum | LGPL-3.0 |
    +----------------------------------+-----------------------------------------+----------+


``Babble``:

..

    +--------------------------------------+--------------------------------------------+--------------+
    |              DEPENDENCY              |                  REPOURL                   |   LICENSE    |
    +--------------------------------------+--------------------------------------------+--------------+
    | github.com/AndreasBriese/bbloom      | https://github.com/AndreasBriese/bbloom    | Other        |
    | github.com/btcsuite/btcd             | https://github.com/btcsuite/btcd           | isc          |
    | github.com/btcsuite/fastsha256       | https://github.com/btcsuite/fastsha256     | Other        |
    | github.com/dgraph-io/badger          | https://github.com/dgraph-io/badger        | Apache-2.0   |
    | github.com/dgryski/go-farm           | https://github.com/dgryski/go-farm         | Other        |
    | github.com/fsnotify/fsnotify         | https://github.com/fsnotify/fsnotify       | bsd-3-clause |
    | github.com/golang/protobuf           | https://github.com/golang/protobuf         | bsd-3-clause |
    | github.com/hashicorp/hcl             | https://github.com/hashicorp/hcl           | MPL-2.0      |
    | github.com/magiconair/properties     | https://github.com/magiconair/properties   | Other        |
    | github.com/mitchellh/mapstructure    | https://github.com/mitchellh/mapstructure  | MIT          |
    | github.com/pelletier/go-toml         | https://github.com/pelletier/go-toml       | MIT          |
    | github.com/pkg/errors                | https://github.com/pkg/errors              | bsd-2-clause |
    | github.com/sirupsen/logrus           | https://github.com/sirupsen/logrus         | MIT          |
    | github.com/spf13/afero               | https://github.com/spf13/afero             | Apache-2.0   |
    | github.com/spf13/cast                | https://github.com/spf13/cast              | MIT          |
    | github.com/spf13/cobra               | https://github.com/spf13/cobra             | Apache-2.0   |
    | github.com/spf13/jwalterweatherman   | https://github.com/spf13/jwalterweatherman | MIT          |
    | github.com/spf13/pflag               | https://github.com/spf13/pflag             | bsd-3-clause |
    | github.com/spf13/viper               | https://github.com/spf13/viper             | MIT          |
    | github.com/ugorji/go                 | https://github.com/ugorji/go               | MIT          |
    | golang.org/x/crypto/ssh/terminal     | https://go.googlesource.com/crypto         |              |
    | golang.org/x/net/internal/timeseries | https://go.googlesource.com/net            |              |
    | golang.org/x/net/trace               | https://go.googlesource.com/net            |              |
    | golang.org/x/sys/unix                | https://go.googlesource.com/sys            |              |
    | golang.org/x/text/transform          | https://go.googlesource.com/text           |              |
    | golang.org/x/text/unicode/norm       | https://go.googlesource.com/text           |              |
    | gopkg.in/yaml.v2                     |                                            |              |
    | github.com/rifflock/lfshook          | https://github.com/rifflock/lfshook        | MIT          |
    +--------------------------------------+--------------------------------------------+--------------+