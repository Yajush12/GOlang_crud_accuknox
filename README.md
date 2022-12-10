# GOlang_crud_accuknox
This is a rest api developed in go language which listen and serve on port 8000. It uses postgreSQL database which uses port 5432.

Some of its functionalities/endpoints include:

1) signup - which request :
  UserIdUserId		string, 
	Email		string, 
	Password	string, 
  
2) login - which request :
  Email		string, 
  Password	string
  
 Its response is :
  Session Id as sid     string
   
3) addNote - which request
  Sid 	string	, 
	Note	string	

4) getNotes - which request :
  Sid		string
 
 Its response is :
  An array containing all notes. UserId, NoteId, Note of the given user as JSON

To test the program use :
1) use "go run main.go" command in terminal
2) use postman on address
    localhost:8000/endpoint_name
use the endpoints mentioned above in place of endpoint_name and mention all the requested stuff as JSON in the body of the request.


Table Schema :

Table 1 - User
UserId		string	primaryKey, 
Email		string 	unique, 
Password	string, 	 
Sid			string	unique, 

Table 2 - Notes
UserId		string	ForeignKey, 
Nid			string	primaryKey, 
Notes		string
