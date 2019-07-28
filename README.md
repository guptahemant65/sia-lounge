# Welcome to SIA Lounge API Documentation

API URL : https://sia-lounge.herokuapp.com

### Following endpoints are live : 

* #### For Client-side App

1. **"/users"**, Methods("GET") - to get all details of all passengers.
2. **"/user"**, Methods("POST") -  to create a new user.
3. **"/user/[:ffn]"**, Methods("GET") - to get details of specific passenger through FFN.
4. **"/user/[:ffn]"**, Methods("PUT") - to modify the details of existing passengers.
5. **"/user/[:ffn]"**, Methods("DELETE") - to delete passenger records.
6. **"/userLogin"**, Methods("POST") - to authenticate the guest login credentials.
7. **"/createLoungeBooking"**, Methods("POST") - to create a new lounge booking.

* #### For Lounge Management App

1. **"/loungeLogin"**, Methods("POST") - to authenticate the lounge login credentials.
2. **"/getLounge/[:ticket_id]"**, Methods("GET") - to get the details of Lounge Booking through Lounge Booking ID/Ticket ID.
3. **"/getLoungeBookings"**, Methods("GET") - to get the all the Lounge Bookings (which are not checked-out yet) on present date. 


#### Sample Data to try out API : 

* FFN(Frequent Flyer Number) : 100254, 140026, 42210002, 42210012
* Lounge ID : 410052, 410042

## Structure of passenger_details Table

| Field Name   |  Data Type                          |  Extras                                |
| ------------ | -------------                       | -------------------------------------- |            
| ffn          |  INT(10)                            | Primary Key, Not Null, Auto Increment  |
| name         |  VARCHAR(45)                        | Not Null                               |
| email        |  VARCHAR(45)                        | Not Null, Unique                       |
| country_code |  VARCHAR(45)                        | Not Null, Unique                       |
| mobile       |  VARCHAR(45)                        | Not Null                               |
| tier_status  |  ENUM(gold,silver,platinum)         | Not Null                               |
| pass         |  VARCHAR(30)                        | Not Null                               |

> Note : Please don't input FFN while using POST Method of "/user" endpoint.  

## Structure of lounge_login Table

| Field Name   |  Data Type                      |  Extras                                |
| ------------ | -------------                   | -------------------------------------- |            
| lounge_id    |  INT                            | Primary Key, Not Null, Auto Increment  |
| pass         |  VARCHAR(45)                    | Not Null                               |

## Structure of lounge_booking Table

| Field Name (description)   |  Data Type                          |  Extras                                |
| ------------               | ------------------                  | -------------------------------------- |            
| ticket_id                  |  VARCHAR(60)                        | Primary Key, Not Null                  |
| ffn                        |  VARCHAR(45)                        | Not Null                               |
| no_of_guests               |  INT(2)                             | Not Null                               |
| guest_names                |  VARCHAR                            | Not Null                               |
| checkin                    |  TIMESTAMP                          |                                        |
| checkout                   |  TIMESTAMP                          |                                        |
| pnr (Flight Ticket ID)     |  VARCHAR                            | Not Null                               |
| status                     | ENUM (CONFIRMED, IN PROGRESS, COMPLETED) | Not Null                          |


> Note : Please don't input ticket_id (as it's auto-generated) while using POST Method of "/createLoungeBooking" endpoint.  



