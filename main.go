package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DiscordUser struct {
	UserID  string `json:"userID"`
	Balance int    `json:"balance"`
}

var (
	ServerId        = flag.String("serverId", "", "server discord id")
	Token           = flag.String("token", "", "token of points api")
	Table           = flag.String("table", "", "Table name")
	Endpoint        = flag.String("endpoint", "http://localhost:8000", "URL to connect database on aws")
	Region          = flag.String("region", "sa-east-1", "Region from aws")
	AccessKeyID     = flag.String("accessKeyId", "i1ie5", "access key aws")
	SecretAccessKey = flag.String("secretAccessKey", "582psh", "secret access key aws")
)

var db *dynamodb.DynamoDB

func init() {
	flag.Parse()
}

func main() {

	api := prepareHttpClient()
	request := prepareRequest()

	res, getErr := api.Do(request)
	if getErr != nil {
		log.Fatal(getErr)
	}

	defer res.Body.Close()

	jsonData := getJSONData(res)

	discordUsers := []DiscordUser{}

	jsonErr := json.Unmarshal([]byte(jsonData), &discordUsers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//WRITE CSV
	file, err := os.Create("points.csv")

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{"userDiscordId", "balance"}

	if err := w.Write(header); err != nil {
		fmt.Println(err)
	}

	for _, discordUser := range discordUsers {
		err = w.Write([]string{discordUser.UserID, strconv.Itoa(discordUser.Balance)})
		if err != nil {
			fmt.Println(err)
		}
	}

	db = CreateLocalClient()
	CreateTableIfNotExists(db, *Table)
	writeDiscordUsers(discordUsers)

	/*
		printVotesOfDb() */
}

func writeDiscordUsers(discordUsers []DiscordUser) {
	for _, discordUser := range discordUsers {
		item, err := dynamodbattribute.MarshalMap(discordUser)
		if err != nil {
			log.Println("Error to write", err)
		}
		_, _ = writeInDynamoDB(db, item, *Table)
	}
}

func prepareHttpClient() http.Client {
	api := http.Client{
		Timeout: time.Second * 30, // Timeout after 2 seconds
	}
	return api
}

func prepareRequest() *http.Request {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.points.city/api/guilds/%s/leaderboard", *ServerId), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *Token))
	return request
}

func getJSONData(res *http.Response) []byte {
	jsonData, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return jsonData
}
