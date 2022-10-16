# Glass
Glass is a tool for listing files in a folder.


## Run with Docker


It is bundled in Docker so that way you can deploy it easily.
It will serve on /data on port 8000.

Just run in the folder you want to serve:
```
docker run --rm -v $PWD:/data -p 8000:8000 decima/glass:1.1
```

### Docker compose

For more convenience you can use docker-compose : 
```yaml
version: "3.7"
services:
  glass:
    image: decima/glass:1.1
    volumes:
      - ./:/data
    ports:
      - 8000:8000
```

## Requirements 
 - golang 16+
 - yarn

## Development

First install project using yarn: 

```
yarn
```

Because for now the project is not well set up, you'll need to do this trick in order to develop.
```
yarn build && go run .
```
You can go on http://localhost:8000 and check the page. If you use a class that was not used, you'll need to restart the command.
I know it's some kind of pain, but it will be fixed in next releases.
