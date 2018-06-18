FROM python:3.6-alpine

WORKDIR /app

COPY requirements.txt /app/requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt

RUN mkdir /usr/local/lib/python3.6/houseparty
COPY __init__.py /usr/local/lib/python3.6/houseparty

# CMD ["python3"]
