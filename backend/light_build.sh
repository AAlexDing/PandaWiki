go build -o panda-wiki-api cmd/api/main.go cmd/api/wire_gen.go
docker cp panda-wiki-api panda-wiki-api:/app/panda-wiki-api
docker restart panda-wiki-api