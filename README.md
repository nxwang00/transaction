Database: PostgreSQL
Starting app: docker-compose build, docker-compose up

?You can run the app and test API endpoints on youu local env without docker.
	-Insall Golang
		https://go.dev/doc/install
	-Install PostgreSQL and pgAdmin
		https://www.postgresql.org/download/

?The project already dockerized so it can be run with docker-compose command like below:

1. Install Docker desktop application.
   	https://docs.docker.com/desktop/install/windows-install/
2. Open the command prompt in the project root.
3. Run "docker-compose build"
4. Run "docker-compose up"
5. Install the Postman to test api endpoints.
6. Send the request to save a transaction record using below API endpoint
   method: POST
   URL: http://localhost:8000/api/v1/trans
   Request Body: {
   				"origin": "desktop-web",
   				"user_id": 1,
   				"amount": "100.00",
   				"op_type": "debit",
   				"registered_at": "2022-10-09 04:05:06"
   			}
7. Send the request to get the list of transactions with filters and paginations
   method: GET,

   -list all transactions
   URL: http://localhost:8000/api/v1/trans

   -list transactions with pagination
   URL: http://localhost:8000/api/v1/trans?page_num=<num>&page_size=<size>
   ex: http://localhost:8000/api/v1/trans?page_num=2&page_size=2

   -list transactions with filters
   URL: http://localhost:8000/api/v1/trans?origin=<origin>&user_id=<id>&amount=<amount>&op_type=<op_type>&registered_at=<registered_at>
   ex: http://localhost:8000/api/v1/trans?origin=desktop-web&user_id=1&amount=100.00&op_type=credit&registered_at=2022-11-03 04:05:06

   -list transactions with pagination and filters
   URL: http://localhost:8000/api/v1/trans?origin=<origin>&user_id=<id>&amount=<amount>&op_type=<op_type>&registered_at=<registered_at>&page_num=<num>&page_size=<size>
