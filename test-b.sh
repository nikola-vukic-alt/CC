#!/bin/bash
echo Trying to borrow a book while not being a member
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina", "Author": "L. N. Tolstoy", "ISBN": "420-3-69-148410-0"}' \
     http://ns.library/borrow
echo
sleep 1

echo Trying to return the book while not being a member
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina"}' \
     http://ns.library/return
echo
sleep 1

echo Registering a new member
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"Name": "Nikola", "Surname": "Vukic", "Address": "Vojvode Supljikca 31", "SSN": "123"}' \
     http://ns.library/register
echo
sleep 1

echo Trying to register a member with the same SSN
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"Name": "Pera", "Surname": "Peric", "Address": "Bulevar Oslobodjenja 105", "SSN": "123"}' \
     http://ns.library/register
echo
sleep 1

echo Trying to return a book while not having borrowed one 
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://ns.library/return
echo
sleep 1

echo Borrowing a book
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina", "Author": "L. N. Tolstoy", "ISBN": "420-3-69-148410-0"}' \
     http://ns.library/borrow
echo
sleep 1

echo Trying to return a book that was not borrowed 
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://ns.library/return
echo
sleep 1

echo Borrowing another book in another library
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace", "Author": "F. M. Dostoevsky", "ISBN": "978-3-16-148410-0"}' \
     http://bg.library/borrow
echo
sleep 1

echo Returning the first book to the first library
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Anna Karenina"}' \
     http://ns.library/return
echo
sleep 1

echo Trying to return the second book to the first library
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://ns.library/return
echo
sleep 1

echo Tryinng to borrow 3 more books while having one borrowed
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Pere Goriot", "Author": "Honore de Balzac", "ISBN": "978-3-16-142411-0"}' \
     http://ns.library/borrow
echo
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary", "Author": "Gustav Flaubert", "ISBN": "696-9-16-148410-0"}' \
     http://nis.library/borrow
echo
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Waiting for Godot", "Author": "Samuel Beckett", "ISBN": "978-3-16-148111-7"}' \
     http://bg.library/borrow
echo
sleep 1

echo Returning the borrowed books
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "War and Peace"}' \
     http://bg.library/return
echo
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Pere Goriot"}' \
     http://ns.library/return
echo
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary"}' \
     http://nis.library/return
echo
sleep 1

echo Trying to return already returned book
sleep 1
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Madame Bovary"}' \
     http://nis.library/return
echo
sleep 1

echo Borrowing the book that could not be borrowed because of limit
sleep 1
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"SSN": "123", "Title": "Waiting for Godot", "Author": "Samuel Beckett", "ISBN": "978-3-16-148111-7"}' \
     http://bg.library/borrow
echo
sleep 1

echo Testing done...
sleep 1
echo Stopping the deployment...
chmod +x k8s/bash/stop-b.sh && ./k8s/bash/stop-b.sh
