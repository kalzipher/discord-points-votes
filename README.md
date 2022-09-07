# Convert JSON from API POINTS CITY DISCORD to CSV

# Create Project
- Run `go mod init github.com/LucioTrucco/discord-points-votes`
- Run `go get .`
- Run `go mod tidy`

## Initial DynamoDBLocal (You need download DynamoDBLocal.jar first)
```shell
docker run -p 8000:8000 amazon/dynamodb-local \
  -jar DynamoDBLocal.jar -sharedDb
```

## HOW TO RUN
- Run `go run . -token "token" -serverId 961074073868308480 -table discord-users-dev -phase phase_1 -endpoint http://localhost:8000 -region sa-east-1 -accessKeyId i1ie5 -secretAccessKey 582psh
2022/08/25 08:55:59 created table=discord-users-dev`
