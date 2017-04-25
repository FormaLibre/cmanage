#!/bin/bash
chmod 777 claroline/app/cache
chmod 777 claroline/app/config
chmod 777 claroline/app/logs
chmod 777 claroline/app/sessions
chmod 777 claroline/web/uploads
chmod 777 claroline/files

docker pull claroline/claroline-docker:prod

source .env
env $(cat .env | grep ^[A-Z] | xargs) bash -c 'docker stack deploy --compose-file docker-compose.yml $PLATFORM_SUBDOMAIN'
echo "Sleeping 30s"
sleep 30

docker exec -it $(docker ps -q --filter name=${PLATFORM_SUBDOMAIN}_claroline) sh -c "cd claroline && php scripts/configure.php"
docker exec -it $(docker ps -q --filter name=${PLATFORM_SUBDOMAIN}_claroline) sh -c "cd claroline && composer install"
docker exec -it $(docker ps -q --filter name=${PLATFORM_SUBDOMAIN}_claroline) sh -c "cd claroline && composer fast-install"
docker exec -it $(docker ps -q --filter name=${PLATFORM_SUBDOMAIN}_claroline) sh -c "sed -i \"/ssl_enabled: false/c\ssl_enabled: true\" claroline/app/config/platform_options.yml"
