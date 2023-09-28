install:
	sudo apt install docker-compose && \
	sudo usermod -aG docker $$USER && \
	sudo service docker restart
rm:
	docker-compose -f docker-compose.yaml stop \
	&& docker-compose -f docker-compose.yaml rm

up:
	docker-compose -f docker-compose.yaml up --force-recreate