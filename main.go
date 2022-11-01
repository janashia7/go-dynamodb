package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamo *dynamodb.DynamoDB

type Person struct {
	Id      int
	Name    string
	Website string
}

const (
	REGION     = "eu-central-1"
	TABLE_NAME = "people"
)

func init() {
	dynamo = connectDynamo()
}

func connectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func createTable() {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

}
func PutItem(person Person) {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
			"Name": {
				S: aws.String(person.Name),
			},
			"Website": {
				S: aws.String(person.Website),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func UpdateItem(person Person) {
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
			"#W": aws.String("Website"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Name": {
				S: aws.String(person.Name),
			},
			":Website": {
				S: aws.String(person.Website),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
		},
		TableName:        aws.String(TABLE_NAME),
		UpdateExpression: aws.String("SET #N = :Name, #W = :Website"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func DeleteItem(id int) {
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func main() {
	// Create Table
	createTable()

	var person Person = Person{
		Id:      1,
		Name:    "George",
		Website: "some.com",
	}

	//  Add Item in table
	PutItem(person)

	// Change Field
	person.Website = "test.com"

	// Update item
	UpdateItem(person)

	// Delete item with partition key Id
	DeleteItem(1)
}
