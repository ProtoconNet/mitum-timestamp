### mitum-timestamp

*mitum-timestamp* is a timestamp contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Installation

Before you build `mitum-timestamp`, make sure to run `docker run` for digest api.

```sh
$ git clone https://github.com/ProtoconNet/mitum-timestamp

$ cd mitum-timestamp

$ go build -o ./mitum-timestamp
```

#### Run

```sh
$ ./mitum-timestamp init --design=<config file> <genesis file>

$ ./mitum-nft run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.
