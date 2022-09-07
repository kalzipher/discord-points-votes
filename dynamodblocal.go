package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateLocalClient() *dynamodb.DynamoDB {

	creds := credentials.NewStaticCredentials(*AccessKeyID, *SecretAccessKey, "")
	awsConfig := &aws.Config{
		Credentials: creds,
	}
	awsConfig.WithRegion(*Region)
	awsConfig.WithEndpoint(*Endpoint)

	s, err := session.NewSession(awsConfig)
	if err != nil {
		panic(err)
	}
	dynamodbconn := dynamodb.New(s)
	return dynamodbconn
}

func CreateTableIfNotExists(d *dynamodb.DynamoDB, tableName string) {
	if tableExists(d, tableName) {
		log.Printf("table=%v already exists\n", tableName)
		d.DeleteTable(&dynamodb.DeleteTableInput{ TableName: &tableName })	
		log.Fatal("Destroy")
		return
	}
	_, err := d.CreateTable(buildCreateTableInput(tableName))
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", tableName)
}

func buildCreateTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("userID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("userID"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}

func tableExists(d *dynamodb.DynamoDB, name string) bool {
	tables, err := d.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {	
		log.Fatal("ListTables failed", err)
	}
	for _, n := range tables.TableNames {
		if *n == name {
			return true
		}
	}
	return false
}

func writeInDynamoDB(d *dynamodb.DynamoDB, item map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.PutItemOutput, error) {
	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	out, err := d.PutItem(itemInput)
	if err != nil {
		fmt.Println(err)
	}
	return out, nil
}
