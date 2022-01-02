# env2dockerenv

A small utility to convert .env files used by the shell to environment
files accepted by the docker --env-file option.

Each line in the normal .env file looks like:

```sh
export VAR="value"
```

A docker env file has just:

```sh
VAR=value
```

To convert, we have to remove any quotes and remove the export statement. As
an additional complication, there exists a convention of using:

```sh
unset VAR
```

This will unset a previous variable setting, effectively commenting out the
value. This utility will prepend the letter 'X' to the environment
variable.

Example usage:

```sh
env2dockerenv -infile .env -outfile .dockerenv
```

Use the -outfile - to preview the output on stdout.
