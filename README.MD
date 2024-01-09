### E2 68-2023 Nikola Vukic

#### A
Start and test the application with 
```shell 
chmod +x test-a.sh && ./test-a.sh
``` 
command (from root directory).\
The script will run *docker compose up*, starting the application, and then send various requests.\
It will test for:
- trying to borrow or return a book by a non-member
- valid registration
- registration with the same social security number
- borrowing book in one but trying to return it in another library
- trying to double borrow
- returning a book that was not borrowed
- trying to double return 
- ...

The results will be displayed in terminal.\
After finishing all the tests it will run *docker compose down*.

#### B
Deploy the application to minikube local cluster with
```shell
chmod +x k8s/bash/deploy-b.sh && ./k8s/bash/deploy-b.sh
```
command (from root directory).\
The script will apply *volumes, claims, deployments, service* and *ingress*.
