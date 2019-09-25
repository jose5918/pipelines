FROM node:9.4.0 as build

ARG COMMIT_HASH
ARG DATE

WORKDIR ./src

COPY . .

WORKDIR ./frontend

# Workaround for ppc64le since phantomjs does not support ppc64le
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        wget --no-verbose -O /tmp/phantomjs-2.1.1-linux-ppc64.tar.bz2 https://github.com/ibmsoe/phantomjs/releases/download/2.1.1/phantomjs-2.1.1-linux-ppc64.tar.bz2 \
        && tar xf /tmp/phantomjs-2.1.1-linux-ppc64.tar.bz2 -C /usr/local/ \
        && ln -s /usr/local/phantomjs-2.1.1-linux-ppc64/bin/phantomjs /usr/bin/phantomjs; \
    fi

RUN npm install && npm run postinstall
RUN npm run build

