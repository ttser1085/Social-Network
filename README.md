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

`curl -v -X POST 'localhost:8094/posts' -H 'Cookie: token= ...' --data '{"title": "My post", "text":"post text"}'`
`curl -v -X GET 'localhost:8094/posts?user_id=aboba'`
`curl -v -X PUT 'localhost:8094/posts' -H 'Cookie: token= ...' --data '{"id": ..., "title": "new My post", "text":"new post text"}'`
`curl -v -X DELETE 'localhost:8094/posts?id=...' -H 'Cookie: token= ...'`
`curl -v -X POST 'localhost:8094/comments' -H 'Cookie: token= ...' --data '{"post_id": ..., "text":"comment text"}'`
`curl -v -X GET 'localhost:8094/comments?post_id=...'`
`curl -v -X DELETE 'localhost:8094/comments?id=...' -H 'Cookie: token= ...'`
