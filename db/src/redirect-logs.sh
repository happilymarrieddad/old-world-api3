#!/bin/sh

make_pipe() {
    pipe=$1
    trap "rm -f $pipe" EXIT

    if [ ! -p $pipe ]; then
        mkfifo $pipe
    fi
}

echo "Using this directory for logs: $NEO4J_server_directories_logs"

# We are piping all logs to stdout so we can see them in Kubernetes 
make_pipe $NEO4J_server_directories_logs/query.log 
make_pipe $NEO4J_server_directories_logs/security.log 
make_pipe $NEO4J_server_directories_logs/debug.log

chown -R neo4j:neo4j $NEO4J_server_directories_logs
cat < /var/lib/neo4j/logs/query.log &
cat < /var/lib/neo4j/logs/security.log &
cat < /var/lib/neo4j/logs/debug.log &

# This is what their (neo4j) docker ENTRYPOINT does
tini -vvv -w -g -- /startup/docker-entrypoint.sh neo4j