## Go-kode

Her finner du instruksjoner for oppsett og gjennomgang av kode.

### Oppsett

<ul>
	<li>Bruk Linux/Mac (jeg bruker WSL2 Ubuntu) og VSCode</li>
	<li><a href="https://go.dev/doc/install" target="_blank">Installer Go v1.20</a></li>
	<li>Installer `@recommended` extensions i VSCode</li>
	<li>
		Go extension i VSCode vil be deg installere noen tools (bl. a. language server), 
		<a href="https://github.com/golang/vscode-go/blob/master/docs/tools.md" target="_blank">gjør det evt manuelt</a>
	</li>
	<li>Github Copilot for assistert læring</li>
</ul>

```bash
wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
go version
source ~/.profile
```

Erstatt eventuelt `linux-amd64` med din arkitektur. Tilgjengelige archives finnes på [go.dev/dl/](https://go.dev/dl/)

### Hello, world!

```bash
go run cmd/hello/main.go
```

### Lab - docker-ctx-analyze

Et CLI som gjør analyse for å finne ut hvor stor docker context blir i bygging av et container image.
Se [cmd/docker-ctx-analyze/main.go](/go/cmd/docker-ctx-analyze/main.go)-filen for instruksjoner.

```bash
go run cmd/docker-ctx-analyze/main.go
```

### Lab - todo-api

Et HTTP API for en todo-applikasjon.
Se [cmd/todo-api/main.go](/go/cmd/todo-api/main.go)-filen for instruksjoner.

```bash
go run cmd/todo-api/main.go
```

### Kjør alle tester

```bash
go test ./...
```
