# go-tiger

CLI for running migrations and generators in go apps.

## Installation

```sh
go install github.com/nrfta/go-tiger/cmd/tiger@latest
```

## Usage

### Configuration

This project uses [go-config](https://github.com/neighborly/go-config) to load
the configuration of your project. We will find the root of your project
(where `go.mod` is saved) and than look at the `config` directory. We will only
read the `JSON` files.

### CLI

Available Commands:

```
  db          Database related operations
  generate    Generate files
  help        Help about any command
  task        Run grift tasks
```

#### Tasks

Run [grift](https://github.com/markbates/grift/) tasks.

#### DB

Available Commands:

| command | description                      |
|---------|----------------------------------|
| create  | Create database                  |
| drop    | Drop database                    |
| migrate | Execute database migrations      |
| reset   | Runs drop, create and migrate up |
| seed    | Runs grift task db:seed          |

Create the development database:

```sh
tiger db create
```

Create the test database:

```sh
ENV=test tiger db create
```

Drop database:
```sh
tiger db drop
```

##### Migrate

We support migration files defined in `db/migrations`. We used [migrate](https://github.com/golang-migrate/migrate) internally.

Available Commands:

| command | description                          |
|---------|--------------------------------------|
| up      | Run migrations in the up direction   |
| down    | Run migrations in the down direction |


#### generate

| command   | description                         |
|-----------|-------------------------------------|
| migration | Generates a migration file          |

## License

This project is licensed under the [MIT License](LICENSE.md).
