# go-webserver-template
Team Support Peru webpage

### Prerequisites
* Go
* Node and npm
* PostgreSQL
* [Air](https://github.com/cosmtrek/air#installation)
* [Templ](https://templ.guide/quick-start/installation)
* inotify-tools

### Install build dependencies
```shell
$ npm install
```

### Initialize the required tables
```shell
$ psql -d <database_name> -U <username> -f ./db/init.sql
```

### .env file example
```
ENV=development
DATABASE_URL=postgres://<username>:<password>@localhost:5432/<database_name>
SESSION_KEY=mysecretkey
PORT=8080
REL=1
SMTP_HOSTNAME=mail.example.com
SMTP_USER=<username>
SMTP_PASS=<password>
```

### Load env variables
```shell
$ export $(grep -v '^#' .env | xargs -d '\n')
```

### Live reload
```shell
$ make live
```
