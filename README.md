# Welcome to SIA Lounge API Documentation

API URL : https://sia-lounge.herokuapp.com

### Following endpoints are live : 

1. **"/users"**, Methods("GET") - to get all details of all passengers.
2. **"/user"**, Methods("POST") -  to create a new user.
3. **"/user/[:ffn]"**, Methods("GET") - to get details of specific passenger through FFN.
4. **"/user/[:ffn]"**, Methods("PUT") - to modify the details of existing passengers.
5. **"/user/[:ffn]"**, Methods("DELETE") - to delete passenger records.
6. **"/loungelogin/[:lounge_id]"**, Methods("GET") - to get login credentials of lounge through Lounge ID.

#### Sample Data to try out API : 

* FFN(Frequent Flyer Number) : 100254, 140026
* Lounge ID : 410002, 410012

## Structure of passenger_details Table

| Field Name   |  Data Type                          |  Extras                                |
| ------------ | -------------                       | -------------------------------------- |            
| ffn          |  INT(10)                            | Primary Key, Not Null, Auto Increment  |
| name         |  VARCHAR(45)                        | Not Null                               |
| email        |  VARCHAR(45)                        | Not Null, Unique                       |
| country_code |  VARCHAR(45)                        | Not Null, Unique                       |
| mobile       |  VARCHAR(45)                        | Not Null                               |
| tier_status  |  ENUM(gold,silver,platinum)         | Not Null                               |
| pass         |  CHAR(76)                           |                                        |

> Note : Please don't input FFN while using POST Method of "/user" endpoint. 

## Structure of lounge_login Table

| Field Name   |  Data Type                          |  Extras                                |
| ------------ | -------------                       | -------------------------------------- |            
| lounge_id    |  INT                                | Primary Key, Not Null, Auto Increment  |
| pass         |  VARCHAR(45)                        | Not Null                               |
