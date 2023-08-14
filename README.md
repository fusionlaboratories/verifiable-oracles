# verifiable-oracles

Implementation of Verifiable Oracles on top of Cosmos SDK.

## devenv

We use [devenv](https://devenv.sh/getting-started/) for specifying the
development environment. It can be either activated using `direnv` or `devenv
shell`.

## Integration Tests with Ganache

CI job is set up to start `ganache` in deterministic mode running on
`127.0.0.1:8545`.  You can manualy start your own instance by using `devenv up`.

1. Start `ganache` in a deterministic mode in a separate terminal window.  It
should bind to localhost (`127.0.0.1`) and its default port (`8545`).  You can
either start the `ganache` process directly

    ```sh
    $ ganache --wallet.deterministic=true
    ganache v7.9.0 (@ganache/cli: 0.10.0, @ganache/core: 0.10.0)
    Starting RPC server
    ...
    ```

    or you can use `devenv up` to start ganache process for you

    ```sh
    $ devenv up
    17:07:58 system    | ganache.1 started (pid=53000)
    17:07:58 ganache.1 | ganache v7.9.0 (@ganache/cli: 0.10.0, @ganache/core: 0.10.0)
    17:07:58 ganache.1 | Starting RPC server
    ...
    ```

2. Run `go test ./... -tags=integration` to run integration tests

    ```sh
    $ go test ./... -tags=integration
    ok      github.com/qredo/verifiable-oracles (cached)
    ok      github.com/qredo/verifiable-oracles/ganache (cached)
    ```
