FROM postgres:15

# Instala pgvector
RUN apt-get update && \
    apt-get install -y postgresql-server-dev-15 build-essential git && \
    git clone --branch v0.5.1 https://github.com/pgvector/pgvector.git && \
    cd pgvector && \
    make && make install && \
    cd .. && rm -rf pgvector && \
    apt-get remove --purge -y build-essential git postgresql-server-dev-15 && \
    apt-get autoremove -y && apt-get clean