# othello

## Run on docker

docker run -it --rm \
    -w "othello" \
    -e "air_wd=othello" \
    -v $(pwd):othello \
    -p 3001:3000 \
    -p 8081:8080 \
    cosmtrek/air
    -c .air.toml
