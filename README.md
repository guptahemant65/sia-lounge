## Welcome to SIA Lounge API Documentation

API URL : https://sia-lounge.herokuapp.com/

Following endpoints are live : 

1. **"/users"**, Methods("GET")
2. **"/user"**, Methods("POST")
3. **"/user/[:ffn]"**, Methods("GET")
4. **"/user/[:ffn]"**, Methods("PUT")
5. **"/user/[:ffn]"**, Methods("DELETE")


Sample FFN(Frequent Flyer Number) to try out API : 100254, 140026

### Structure of passenger_details Table

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
