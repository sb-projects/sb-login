down:
	docker compose -f docker-compose.yaml down
# docker image rm sb-login-login-service -f

up:
	docker compose -f docker-compose.yaml up -d

service_up:
	docker compose -f docker-compose.yaml up -d --no-deps --build login-service

# Remove all <none> images
prune:
	docker image prune

kafka_down:
	docker compose -f docker-compose.kafka.yaml down

kafka_up:
	docker compose -f docker-compose.kafka.yaml up -d

kafka_topic:
	docker compose up -d --no-deps --build kafka-helper

up_all:
	docker compose -f docker-compose.kafka.yaml up -d
	docker compose -f docker-compose.multi.yaml up -d