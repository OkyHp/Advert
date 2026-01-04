#### Собираем builder image (linux/amd64)
```sh
docker build --platform linux/amd64 -t steamrt-go-builder -f Dockerfile .
```

#### Запускаем сборку и получаем бинарник прямо в проект
```sh
docker run --rm \
  --platform linux/amd64 \
  --user "$(id -u):$(id -g)" \
  -e BIN_NAME=Advert.so \
  -e BUILD_DEBUG=false \
  -v "$(pwd):/app" \
  # Add volume for dependency
  -v "$(pwd)/../plg_utils:/plg_utils" \
  -w /app \
  steamrt-go-builder
```



