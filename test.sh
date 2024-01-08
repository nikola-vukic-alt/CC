#!/bin/bash

make_request() {
  local url=$1
  local method=$2
  local data=$3
  local response=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -d "$data" $url)
  local body=$(echo "$response" | head -n -1)  # Exclude the last line (status code)
  local status=$(echo "$response" | tail -n 1)  # Get the last line (status code)
  echo "Status Code: $status"
  echo "Response Body: $body"
}

docker compose up -d
sleep 3

echo Trying to borrow a book while not being a member
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina", "Author": "L. N. Tolstoy", "ISBN": "420-3-69-148410-0"}' \
     http://localhost:8081/borrow
sleep 1

echo Trying to return the book while not being a member
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina"}' \
     http://localhost:8081/return
sleep 1

echo Registering a new member
sleep 1
url="http://localhost:8080/register"
method="POST"
body='{"Name": "Nikola", "Surname": "Vukic", "Address": "Vojvode Supljikca 31", "SSN": "123"}'
make_request $url $method "$body"
# curl -X POST \
#      -H "Content-Type: application/json" \
#      -d '{"Name": "Nikola", "Surname": "Vukic", "Address": "Vojvode Supljikca 31", "SSN": "123"}' \
#      http://localhost:8080/register
sleep 1

echo Trying to register a member with the same SSN
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"Name": "Pera", "Surname": "Peric", "Address": "Bulevar Oslobodjenja 105", "SSN": "123"}' \
     http://localhost:8080/register
sleep 1

echo Trying to return a book while not having borrowed one 
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://localhost:8081/return
sleep 1

echo Borrowing a book
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina", "Author": "L. N. Tolstoy", "ISBN": "420-3-69-148410-0"}' \
     http://localhost:8081/borrow
sleep 1

echo
echo Trying to return a book that was not borrowed 
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://localhost:8081/return
sleep 1

echo Borrowing another book in another library
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace", "Author": "F. M. Dostoevsky", "ISBN": "978-3-16-148410-0"}' \
     http://localhost:8082/borrow
sleep 1

echo
echo Returning the first book to the first library
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina"}' \
     http://localhost:8081/return
sleep 1

echo Trying to return the second book to the first library
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://localhost:8081/return
sleep 1

echo Tryinng to borrow 3 more books while having one borrowed
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Pere Goriot", "Author": "Honore de Balzac", "ISBN": "978-3-16-142411-0"}' \
     http://localhost:8081/borrow
echo
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary", "Author": "Gustav Flaubert", "ISBN": "696-9-16-148410-0"}' \
     http://localhost:8083/borrow
echo
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Waiting for Godot", "Author": "Samuel Beckett", "ISBN": "978-3-16-148111-7"}' \
     http://localhost:8082/borrow
sleep 1

echo Returning the borrowed books
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://localhost:8082/return
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Pere Goriot"}' \
     http://localhost:8081/return
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary"}' \
     http://localhost:8083/return
sleep 1

echo Trying to return already returned book
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary"}' \
     http://localhost:8083/return
sleep 1

echo Borrowing the book that could not be borrowed because of limit
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Waiting for Godot", "Author": "Samuel Beckett", "ISBN": "978-3-16-148111-7"}' \
     http://localhost:8082/borrow
sleep 1

echo
echo Testing done...
sleep 1
docker compose down 
sleep 3
