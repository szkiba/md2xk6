# md2xk6

**Use markdown file as xk6 build manifest**

> [!WARNING]
> This is a temporary development repository specifically designed for the `grafana/H-14-my-k6-template` repository. It is not intended for general use and may be removed or substantially modified without prior notice.

**md2xk6** generates `--with` flags for `xk6 build` from the appropriate list in the markdown file (`default README.md`).

You can specify extensions using any of these formats:

- **Markdown links**: `[extension-name](https://github.com/user/repo)`
- **Auto-links**: `<https://github.com/user/repo>`
- **Plain URLs**: `https://github.com/user/repo`
- **Versioned releases**: `[extension-name](https://github.com/user/repo/releases/tag/v1.0.0)`

Configuration Rules:

- Include exactly one extension per list item
- Only the first valid list in this README will be processed
- Content before and after the extension list is ignored
- There can be no text before or after the plain URL

## Example

Our team uses the following k6 extensions:

- [grafana/xk6-faker](https://github.com/grafana/xk6-faker) for generating random data.
- <https://gitlab.com/szkiba/xk6-banner> for generating funny ASCII banners
- [xk6-sql v1.0.0](https://github.com/grafana/xk6-sql/releases/tag/v1.0.0) for database management
- <https://github.com/grafana/xk6-sql-driver-ramsql>
- https://github.com/grafana/xk6-exec

## Contributing

We use [Development Containers](https://containers.dev/) to provide a reproducible development environment. We recommend that you do the same. In this way, it is guaranteed that the appropriate version of the tools required for development will be available.

### Tasks

The usual contributor tasks can be performed using GNU make. The `Makefile` defines a target for each task. To execute a task, the name of the task must be specified as an argument to the make command.

```bash
make taskname
```

Help on the available targets and their descriptions can be obtained by issuing the `make` command without any arguments.

```bash
make
```

More detailed help can be obtained for individual tasks using the [cdo](https://github.com/szkiba/cdo) command:

```bash
cdo taskname --help
```

Authoring the Makefile

The `Makefile` is generated from the task list defined in the `CONTRIBUTING.md` file using the [cdo](https://github.com/szkiba/cdo) tool. If a contribution has been made to the task list, the `Makefile` must be regenerated using the [makefile] target.

```bash
make makefile
```

#### lint - Run the linter

The [golangci-lint] tool is used for static analysis of the source code. It is advisable to run it before committing the changes.

```bash
golangci-lint run ./...
```

[lint]: <#lint---run-the-linter>
[golangci-lint]: https://github.com/golangci/golangci-lint

#### security - Run security and vulnerability checks

The [gosec] tool is used for security checks. The [govulncheck] tool is used to check the vulnerability of dependencies.

```bash
gosec -quiet ./...
govulncheck ./...
```

[gosec]: https://github.com/securego/gosec
[govulncheck]: https://github.com/golang/vuln
[security]: <#security---run-security-and-vulnerability-checks>

#### test - Run the tests

The `go test` command is used to run the tests and generate the coverage report.

```bash
go test -count 1 -race -timeout 2m ./...
```

[test]: <#test---run-the-tests>

### build - Build executable

The [goreleaser] tool is used to build the executable.

```bash
goreleaser build --clean --snapshot --single-target
```

[build]: <#build---build-executable>

### clean - Clean the working directory

Delete the work files created in the work directory (also included in .gitignore).

```bash
rm -rf ./md2xk6 ./md2xk6.exe ./k6 ./k6.exe ./build ./dist
```

[clean]: #clean---clean-the-working-directory

### all - Clean build

Performs the most important tasks. It can be used to check whether the CI workflow will run successfully.

Requires
: [clean], [format], [security], [lint], [test], [build], [makefile]

### format - Format the go source codes

```bash
go fmt ./...
```

[format]: #format---format-the-go-source-codes

### makefile - Generate the Makefile

```bash
cdo --makefile Makefile
```

[makefile]: <#makefile---generate-the-makefile>
