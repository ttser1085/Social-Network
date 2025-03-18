# Social Network
Багдасарян А. А. БПМИ225

### Auth service examples:

- signup:
    `curl -v -X POST 'localhost:8091/signup' --data '{"id": "aboba", "name": "boba", "password": "pass123"}'`

- login:
    `curl -v -X POST 'localhost:8091/login' --data '{"id": "aboba", "password": "pass123"}'`

- whoami:
    `curl -v -X GET 'localhost:8091/whoami' -H 'Cookie: jwt= ... '`

- update:
    `curl -v -X GET 'localhost:8091/update' --data '{"name": "boba2"}'  -H 'Cookie: jwt= ...'`
