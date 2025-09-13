FROM python:3.9-alpine

WORKDIR /practice-security/tests

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY *.py .

ENTRYPOINT pytest -vs --log-level=DEBUG
