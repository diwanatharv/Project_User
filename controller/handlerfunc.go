package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

var Id = 1
var v = validator.New()

type User struct {
	Id           int    `json:"id" bson:"id" validate:"omitempty"`
	FirstName    string `json:"firstName" bson:"firstName" validate:"required,gte=2"`
	LastName     string `json:"lastName" bson:"lastName"`
	Email        string `json:"email" bson:"email"validate:"required_with=Id,email"`
	BusinessType string `json:"businessType" bson:"businessType"`
	PhoneNo      string `json:"phoneNo" bson:"phoneNo" validate:"number"`
	CompanyName  string `json:"companyName" bson:"companyName" `
	Country      string `json:"country" bson:"country"validate:"required"`
}

func Iserror(err error) bool {
	if err != nil {
		return true
	}
	return false
}
func CreateUser(e echo.Context) error {
	var instance User
	err1 := e.Bind(&instance) // this  binds the instance of the user struct from the users data
	if Iserror(err1) {
		return e.JSON(http.StatusNotFound, "error in intializing ")
	}
	if err2 := v.Struct(&instance); Iserror(err2) { // this used for validation
		return e.String(http.StatusForbidden, "validation failed try again")
	}
	// for giving the unique id counting the no of documents inside the collection and inserting its
	//lenght as id
	ans, _ := User_Collection.EstimatedDocumentCount(context.Background()) //giving unique id
	instance.Id = (int)(ans) + 1
	res, err3 := User_Collection.InsertOne(context.Background(), instance)
	fmt.Println(res)
	if Iserror(err3) {
		fmt.Println("problem in insertion")
		return e.String(http.StatusForbidden, "problem in inserting elements to the mongodb")
	}
	return e.JSON(http.StatusOK, res)
}

func GetUser(e echo.Context) error {
	//getting the value what value is the user trying to search on
	UserId, err1 := strconv.Atoi(e.QueryParam("userId"))
	if Iserror(err1) {
		return e.JSON(http.StatusNotFound, "error in params ")
	}
	//creating an instance of the user finding in db
	// if we are getting the value correponding to that query we are storing it in
	// the finduser and returning it if no error
	val, err3 := Rdb.Get(context.Background(), e.QueryParam("userId")).Result()
	if err3 == redis.Nil { //means the val is not present we have to add it
		fmt.Println("value is not there")
		var findUser User
		err2 := User_Collection.FindOne(context.Background(), bson.M{"id": UserId}).Decode(&findUser)
		if Iserror(err2) {
			return e.JSON(http.StatusNotFound, "error in finding the relevant id")
		}
		find_New_User, _ := json.Marshal(findUser) //converts in json before setting th value of the users

		key := strconv.Itoa(findUser.Id) //which uniquely identifies it
		err4 := Rdb.SetNX(context.Background(), key, find_New_User, redis.KeepTTL).Err()

		if Iserror(err4) {
			return e.JSON(http.StatusNotFound, "error in setting the value to the reddis")
		}
		return e.JSON(http.StatusOK, findUser)
	} else {
		var j User
		err := json.Unmarshal([]byte(val), &j)
		if Iserror(err) {
			return e.JSON(http.StatusNotFound, "error in unmarshal")
		}
		return e.JSON(http.StatusOK, j)
	}

}
func getAllUser(c echo.Context) error {
	var allUser []User
	var oneUser User
	// this makes an unordered map with id and marks true
	elementFilter := bson.M{
		"id": bson.M{"$exists": true},
	}
	findElementRes, err := User_Collection.Find(context.Background(), elementFilter)
	// returns a cursor struct
	if err != nil {
		return err
	}
	for findElementRes.Next(context.Background()) {
		// decoding the value and putting it in one user
		err := findElementRes.Decode(&oneUser)
		if err != nil {
			fmt.Println(err)
		}
		// inserting it in the alluser slice
		allUser = append(allUser, oneUser)
	}
	// returning the slice of the array
	return c.JSON(http.StatusOK, allUser)
}
func UpdateUser(c echo.Context) error {
	//leadId
	UserId, err1 := strconv.Atoi(c.QueryParam("userId"))

	if Iserror(err1) {
		fmt.Println("1 error")
		return c.JSON(http.StatusForbidden, "")
	}
	//isPresent
	var findUser User
	err2 := User_Collection.FindOne(context.Background(), bson.M{"id": UserId}).Decode(&findUser)

	if Iserror(err2) {
		return c.JSON(http.StatusForbidden, "Not present .first create the user")
	}
	var reqBody User
	err3 := c.Bind(&reqBody) // changes saved
	if Iserror(err3) {
		fmt.Println("3 error")
		return c.JSON(http.StatusForbidden, "Bind problem")
	}

	err4 := v.Struct(reqBody)
	if Iserror(err4) {
		fmt.Println("4 error")
		return c.JSON(http.StatusForbidden, "validation problem")
	}
	updateFieldCon := bson.M{"$set": bson.M{"firstname": reqBody.FirstName, "lastname": reqBody.LastName, "email": reqBody.Email, "businessType": reqBody.BusinessType, "phoneNo": reqBody.PhoneNo, "companyName": reqBody.CompanyName, "country": reqBody.Country}}
	updateFileRes, err := User_Collection.UpdateOne(context.Background(), bson.M{"id": findUser.Id}, updateFieldCon)
	reqBody.Id = UserId
	if Iserror(err) {
		fmt.Println("4 error")
		return c.JSON(http.StatusForbidden, "validation problem")
	}
	reqBody2, _ := json.Marshal(reqBody)
	fmt.Println(reqBody)
	err7 := Rdb.Set(context.Background(), strconv.Itoa(findUser.Id), reqBody2, 0).Err()
	if err7 == redis.Nil {
		fmt.Println("not able to set the values in redis")
	}
	return c.JSON(http.StatusOK, updateFileRes)
}
