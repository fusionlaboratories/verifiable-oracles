# verifiable-oracles

Implementation of Verifiable Oracles on top of Cosmos SDK.

## devenv

We use [devenv](https://devenv.sh/getting-started/) for specifying the
development environment. It can be either activated using `direnv` or `devenv
shell`.

## Ganache

This is a Ganache client written in Go.  We will be using it to simulate
transactions in order to test the verifiable oracle.

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

2. Run `go run ./ganache` to run the client.

    ```sh
    $ go run ./ganache
    Successfully connected to ganache on http://127.0.0.1:8545
    Current block number is 0
    Account 0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1 has 1000 ETH
    Account 0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0 has 1000 ETH
    ```
