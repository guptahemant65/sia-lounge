# Welcome to SIA Lounge API Documentation

API URL : https://sia-lounge.herokuapp.com

### Following endpoints are live : 

* #### For Client-side App

1. **"/users"**, Methods("GET") - to get all details of all passengers.
2. **"/user"**, Methods("POST") -  to create a new user. </br>
   **Params Required** : name, email, country_code, mobile, tier_status, pass </br>
3. **"/user/[:ffn]"**, Methods("GET") - to get details of specific passenger through FFN.
4. **"/user/[:ffn]"**, Methods("PUT") - to modify the details of existing passengers.
5. **"/user/[:ffn]"**, Methods("DELETE") - to delete passenger records.
6. **"/userLogin"**, Methods("POST") - to authenticate the guest login credentials. </br>
   **Params Required** : ffn, pass (Password) </br>
7. **"/createLoungeBooking"**, Methods("POST") - to create a new lounge booking. </br>
   **Params Required** : ffn, lounge_id, no_of_guests, guest_name, checkin, pnr, payment_method, amount_paid </br>
8. **"/getpnr/[:ffn]"**, Methods("GET") - to get all the upcoming fight bookings for an entered FFN.

* #### For Lounge Management App

1. **"/loungeLogin"**, Methods("POST") - to authenticate the lounge login credentials. </br>
   **Params Required** : lounge_id, pass (password) </br>
2. **"/getLoungeBooking/[:ticket_id]"**, Methods("GET") - to get the details of Lounge Booking through Lounge Booking ID/Ticket ID.
3. **"/getUpcomingLoungeBookings/[:lounge_id]"**, Methods("GET") - to get the all the upcoming Lounge Bookings.
4. **"/getCurrentLoungeBookings/[:lounge_id]"**, Methods("GET") - to get the all the Lounge Bookings (which are checked-in) on present date. 
5. **"/getLoungeDetails/[:lounge_id]"**, Methods("GET") - to get the Lounge Details with Lounge_ID.
6. **"/checkin"**, Methods("POST") </br>
   **Params Required** : ticket_id </br>
7. **"/checkout"**, Methods("POST") </br>
   **Params Required** : ticket_id </br>
8. **"/getLoungeBookings/[:ffn]"**, Methods("GET") - to get the details of all Lounge Bookings for a particular FFN.
9. **"/getLoungeDetails"**, Methods("GET") - to get all Lounge Details.

* #### Payment Related

1. **"/cardCheck"**, Methods("POST") -  to get the number of complimentary lounges left and list of available lounges for that card. </br>
   **Params Required** : card_number </br>
   

#### Sample Data to try out API : 

* FFN(Frequent Flyer Number) : 100254, 140026, 42210002, 42210012
* Lounge ID : 410052, 410042
* Card Number : 4577044405321376, 5589831300967730

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

| Field Name (description)   |  Data Type                               |  Extras                                |
| -------------------------- | ---------------------------------------  | -------------------------------------- |            
| ticket_id                  |  VARCHAR(60)                             | Primary Key, Not Null                  |
| ffn                        |  VARCHAR(45)                             | Not Null                               |
|lounge_id                   |                                          | Not Null                               |
|lounge_name                 |                                          | Not Null                               |
|lounge_address              |                                          | Not Null                               |
| no_of_guests               |  INT(2)                                  | Not Null                               |
| guest_names                |  VARCHAR                                 | Not Null                               |
| checkin                    |  TIMESTAMP                               |                                        |
| checkout                   |  TIMESTAMP                               |                                        |
| pnr (Flight Ticket ID)     |  VARCHAR                                 | Not Null                               |
| status                     | ENUM (CONFIRMED, IN PROGRESS, COMPLETED) | Not Null                               |
| payment_method             | VARCHAR                                  | Not Null                               |
| amount_paid                | INT                                      | Not Null                               |


> Note : Please don't input ticket_id (as it's auto-generated) while using POST Method of "/createLoungeBooking" endpoint.  

## Structure of lounge_details Table

| Field Name (description)   |  Data Type                          |  Extras                                |
| ------------               | ------------------                  | -------------------------------------- |  
| lounge_id                  | INT(11)                             | Not Null,Primary Key                   |
| lounge_name                | VARCHAR(45)                         | Not Null                               |
| total_capacity             | INT(3)                              | Not Null                               |
| amenities                  | VARCHAR(120)                        | Not Null                               |
| price                      | INT(4)                              | Not Null                               |
|accepted_cards              | VARCHAR(150)                        | Not Null                               |
| private_room_capacity      | INT(2)                              | Not Null                               |
| sofa_capacity              | INT(3)                              | Not Null                               |
| location                   | VARCHAR(45)                         | Not Null                               |

## Structure of card_details Table

| Field Name (description)   |  Data Type                          |  Extras                                |
| ------------               | ------------------                  | -------------------------------------- | 
| card_number                |                                     | PRIMARY KEY, NOT NULL                  |
| lounge_left                |                                     | NOT NULL                               |
| available_lounge           |                                     | NOT NULL                               | 



