# GoDBLedger-Web
Webserver GUI for GoDBLedger

Still very much in Alpha

## How to Use

Build the server using
```
make
````
then run the webserver on the same host as a GoDBLedger server.
```
make run
```

This will open a webserver over port :3000 that you can navigate to with your webbrowser.

## Docker 

Godbledger comes with a `docker-compose.yml` file and some make targets to help build the `godbledger` server and the `godbledger-web` server into a docker container and launch it with a mysql backend, configuring both to store state inside the host's default DATA_DIR so that state persists by default across restarts of the containers.

1. Build the container image

    ```
    make docker-build
    ```

    This builds `godbledger-web` into an alpine-based container image tagged locally as `godbledger-web`

1. Start mysql and godbledger server in docker

    ```
    make docker-start
    ```

    This invokes `docker-compose up` with a few env vars set for default configuration.

    There are some env vars which can adjust configuration but by default:
    - `mysql` is running and reachable through docker at `localhost:3306`
      - a `ledger` database has been created along with the following local user account:
        - username: `godbledger`
        - password: `password`
    - `godbledger` server is available through docker at `localhost:50051` and configured to use that mysql service as a backend
    - `godbledger-web` server is available through docker at `localhost:3000` and configured to use that mysql service as a backend
    - CLI tools running on your local host machine can connect with the following values in your local `config.toml` file:

        ```toml
        Host = "127.0.0.1"
        RPCPort = "50051"
        DatabaseType = "mysql"
        DatabaseLocation = "godbledger:password@tcp(localhost:3306)/ledger?charset=utf8mb4,utf8
        ```

1. Stop mysql and godbledger server

   In the terminal where you ran `make docker-start` you can use `ctrl-c` to gracefully shut down both containers.

   From another terminal you can run `make docker-stop` which invokes `docker-compose down`

   Because the database and docker `config.toml` files are stored on your host machine, you can safely stop and start both apps without losing data.

   **NOTE** you may risk losing data if you type `ctrl-c` twice and force an early shutdown of mysql

## Authentication

User: test@godbledger.com
Pass: password

### Roadmap
Discussion can be found on this github issue:
https://github.com/darcys22/godbledger/issues/169


### Misc
https://github.com/mattrobenolt/grafana/tree/ae0cbdff73297bd8ab24b6a9bbfdcb0cd1439218

https://select2.org/
https://github.com/select2/select2-bootstrap-theme
