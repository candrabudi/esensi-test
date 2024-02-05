
# Esensi Test

Procedures for project initialization


## Installation

The first step is to clone the GitHub project

```bash
https://github.com/candrabudi/esensi-test
cd esensi-test
```

copy the env first

```bash
cp .env.example .env
```

generate JWT_KEY using bcrypt generator, then copy the bcrypt results and add JWT_KEY in .env

```bash
https://bcrypt--generator-com.webpkgcache.com/doc/-/s/bcrypt-generator.com/
JWT_KEY="example"
```

After the setup above, don't forget to change the database settings in the env, and create a MySQL database

To migrate tables, you can use the command

```bash
go run main.go -m="migrate"
```

then to run the project you can use

```bash
go run main.go
```