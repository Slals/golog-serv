# Ultra Extra Simple Remote Logger Service

I needed a simple remote logger service for tracing what's going on in mobile applications and view it in a html page. So I wrote it in Go.

Logs data is saved in files, no query langage used, only appending.

## Installation 

Simply run `docker build --tag golog-serv .` and `docker run -d -p 127.0.0.1:3333:3333 --name golog-serv golog-serv` 

## Env vars

You can edit them from Dockerfile.

DEBUG_PATH is the path used to save debug files.

PAGE_TITLE is the page title of the logger page.

## Usage

```
Request
PUT /logs
Content-Type: application/json
{
    key: "key_message",
    message: "your message",
    level: "log level"
}

Response
Status: 204 NoContent
```

Available levels are listed bellow:

- "trace": Used to keep track of normal processes
- "debug": Same as trace but only used for development environment
- "info": Used to keep track of scheduled operations
- "notice": Used to track noticable event from production environment
- "warn": Used to track events that could lead to an error
- "error": Used to track errors which doesn't kill the client process from develpment environment and / or production environment
- "fatal": Used to track fatal errors which kill the client process from development environment and / or production environment

```
Request
GET /logs
Accept: text/html

Response
index.html (the page which shows the data)
```

If you want to protect `GET /logs`, there is no authentication yet, use the basic authentication of your http server, for instance to protect it with NginX I use this instruction :

```
limit_except PUT {
    auth_basic "Login required";

    # .htpwasswd created by using apacheutils
    auth_basic_user_file /secret/path/.htpasswd;
}
```

## TODO and contribution

There is some work to do which I would probably do in the future, here is a list if you want to kindly contribute :

- [ ] Make the index.html prettier with filters.
- [ ] Write logs on multiple files for scalability.
- [ ] Adds HMAC based authentication. 

_Any suggestion? Post an issue! Thanks :-)_

