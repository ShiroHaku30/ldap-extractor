# Instructions on Installation and Build

## Installation

### Linux

1. Install go version 1.27 on the machine

> https://go.dev/dl/

2. Ensure the below dependencies are installed

- make

### Windows

1. Ensure the below dependencies are installed

- cygwin or mingw
- make

2. Install go version 1.27 on the machine

> https://go.dev/dl/

## Build

1. Run the make build command

```
make build
```

---


# Command Usage

### Dump

For Dumping LDAP objects based on configured config.yaml ( testing purposes )

```bash
ldap-extractor dump -c <config.yaml>
```

### Extract

For dumping, filtering, and generating diff between between latest filtered and previous filtered.

```bash
ldap-extractor extract -c <config.yaml>
```

### Filter

For Filtering out LDAP users based on the dumped ldif file ( testing purposes )

```bash
ldap-extractor filter -i <ldap_dump.ldif> -o <filtered.json>
```

### Test Connection

For testing the ldap connection to check if auth succeeds based on provided credentials

```bash
ldap-extractor test -c <config.yaml>
```