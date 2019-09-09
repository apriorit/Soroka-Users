#!/usr/bin/env sh

echo "Adding keys to k/v storage..."

# add connection string to consul
DB_CONNECTION="stub"
curl -X PUT --data-binary "$DB_CONNECTION" http://$CONSUL_HOST_ADDR/v1/kv/$USERS_DB_KEY

chmod +x /bin/userssvc

/bin/userssvc -consul.address $CONSUL_HOST_ADDR