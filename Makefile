build:
	docker build -t mtgdrop-server -f ./server/prod.dockerfile ./server

run:
	docker run -p 4000:4000 mtgdrop-server

