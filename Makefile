dev-backend:
	cd backend && go run main.go
dev-ui:
	cd ui && yarn start

IMAGE=nicky/food-backend

docker-build:
	docker build -t $(IMAGE) .
docker-run:
	docker run -p 8080:8080 -e "DB_USERNAME=root" -e "DB_PASSWORD=6NzuCvnlO3cwOFvs" -e "DB_DATABASE=food" -e "DB_DATABASE_TEST=food2" -e "DB_HOST=35.230.30.35" -e "PORT=8080" $(IMAGE) 
docker-dev: docker-build docker-run

docker-push: docker-build
	docker push $(IMAGE):latest