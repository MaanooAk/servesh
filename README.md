# servesh

Web server that passes the requests to a script.

**servesh** was created to easily setup a server with a shell script or just run something from a browser remotely.

```
servesh some-program
```

```
http://localhost:8080/
 V
some-program
```

```
http://localhost:8080/hello?name=master
 V
some-program hello --name master
```

Example where json file outputed formated after a simulated delay:

```
servesh 'sleep 2; cat ./samples/{}.json | jd'
```

Example where any request will return a screenshot of the of the entire screen (powered by imagemagick):

```
servesh 'import -window root png:-'
```

## Usage

**servesh** takes the first non option argument at handles it as the shell command.

Each http request is converted to command line arguments. Every instace of `{}` in the shell command is replaced with the generated command line arguments. If no `{}` is contained in the shell command then ` {}` is appended at the end. To prevent this a `;` can be added at the end of the shell command.

The path of the http request is converted to the first arguments. Every key-value pair of the url's query are converted to arguments with values (eg. `...?action=view&id=10` to `--action view --id 10`). The prefix `--` of the key can be change.

The shell command is executed by passing it to `sh -c`.

### Options

- `--port INT`: set the server port (default: 8080)
- `--prefix TEXT`: set the prefix for the url query parameters keys (default: `--`)
- `--no-prefix`: disable the prefix for the url query parameters keys (equivalent to `--prefix ''`)
- `--icos`: allow requests for *.ico files. Ico file requests are discarded by default for better handling of the `favicon.ico` request that browsers send.


## Install

```
git clone https://github.com/MaanooAk/servesh.git
cd servesh
go install
```

