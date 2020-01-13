# go-tiger

CLI for running migrations and generators in go apps.

## Installation

```sh
go install github.com/nrfta/go-tiger/cmd
```

## Usage

### Configuration

This project uses [go-config](https://github.com/neighborly/go-config) to load
the configuration of your project. We will find the root of your project
(where `go.mod` is saved) and than look at the `config` directory. We will only
read the `JSON` files.

### CLI

#### DB

Available Commands:

| command | description                 |
|---------|-----------------------------|
| create  | Create database             |
| drop    | Drop database               |
| migrate | Execute database migrations |

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
