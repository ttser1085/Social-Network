import pytest
import requests
import string
import random

import logging

LOGGER = logging.getLogger(__name__)

@pytest.fixture
def signup():
    return 'http://auth:8091/signup'

@pytest.fixture
def posts():
    return 'http://grpc-gateway:8094/posts'

@pytest.fixture
def comments():
    return 'http://grpc-gateway:8094/comments'

def random_string(length=8):
    chars = string.ascii_letters
    return ''.join(random.choices(chars, k=length))

def test_endtoend(signup, posts, comments):
    user = random_string()
    LOGGER.info(user)

    session = requests.Session()

    req = requests.Request(method='POST', url=signup, data='{"id": "'+user+'", "name": "aboba", "password": "pass123"}', ).prepare()
    resp = session.send(req)
    LOGGER.debug(resp.content)
    assert resp.status_code == 200

    token = resp.cookies['jwt']
    LOGGER.info(token)

    req = requests.Request(method='POST', url=posts, data='{"title": "post title1", "text":"post text1"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='POST', url=posts, data='{"title": "post title2", "text":"post text1"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='GET', url=posts+'?user_id='+user).prepare()
    resp = session.send(req)
    assert resp.status_code == 200
    assert resp.content.decode().count('id') == 2

    content = str(resp.content)
    post_id = content[content.find('"id":"')+6:content.find('","')]
    LOGGER.info(post_id)

    req = requests.Request(method='PUT', url=posts, data='{"id":"'+post_id+'", "title": "new post title", "text":"new post text"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='GET', url=posts+'?user_id='+user).prepare()
    resp = session.send(req)
    assert resp.status_code == 200
    assert str(resp.content).count('new post') == 2

    req = requests.Request(method='DELETE', url=posts+'?id='+post_id, cookies={'token': token}).prepare()
    resp = session.send(req)
    LOGGER.debug(resp.content)
    assert resp.status_code == 200

    req = requests.Request(method='GET', url=posts+'?user_id='+user).prepare()
    resp = session.send(req)
    assert resp.status_code == 200
    assert resp.content.decode().count('id') == 1

    content = str(resp.content)
    post_id = content[content.find('"id":"')+6:content.find('","')]

    req = requests.Request(method='POST', url=comments, data='{"post_id": "'+post_id+'", "text":"ABOBA"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='POST', url=comments, data='{"post_id": "'+post_id+'", "text":"ABOBA"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='GET', url=comments+'?post_id='+post_id).prepare()
    resp = session.send(req)
    LOGGER.debug(resp.content)
    assert resp.status_code == 200
    assert str(resp.content).count('ABOBA') == 2

    content = str(resp.content)
    comm_id = content[content.find('"id":"')+6:content.find('","')]

    req = requests.Request(method='PUT', url=comments, data='{"id": "'+comm_id+'", "text":"BOBA"}', cookies={'token': token}).prepare()
    resp = session.send(req)
    assert resp.status_code == 200

    req = requests.Request(method='GET', url=comments+'?post_id='+post_id).prepare()
    resp = session.send(req)
    LOGGER.debug(resp.content)
    assert resp.status_code == 200
    assert str(resp.content).count('ABOBA') == 1
