# Go API

This application is a simple REST API which provides one serviceto run a provided dockerfile.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

To run this project you'll need a Docker and/or Golang working environment.<br>
Install Docker compose : https://docs.docker.com/compose/install/<br>

### Installing

Step-by-step installation guide
Clone repo following go convention:
```
 cd ~/go/src/github.com/alisstaki
 git clone https://github.com/alisstaki/interview-exercise.git
 cd interview-exercise
```
Start app
```
make start
```

The app should be up and running and see something like that in your terminal:
```
API ready to receive content...
```

## Dependencies

* Third parties libraries used : https://github.com/gorilla/mux to manage services routing and https://github.com/moby/buildkit to validate dockerfile

## Running the tests

For time optimization, I've decided to not unittest this api

## Consuming the API

You can upload your dockerfile using a POST to http://localhost:8080/content :
```
curl --request POST \
  --url http://localhost:8080/content \
  --header 'Content-Type: multipart/form-data; boundary=---011000010111000001101001' \
  --form 'file=FROM ubuntu:16.04

# train machine learning model

# save performances

CMD echo '\''{"perf":0.99}'\'' > /data/perf.json'
```

## To go further

With given limited time, the focus was made to have a running project and a validity check on the dockerfile.
To continue, I would focus on security and only allow a certain set of command in the provided dockerfile, and run it in a more secure environment.
Also jwt token authentication would be added to allow only known client to upload a dockerfile.

At last, additionnal time could be spent to have a cleaner code and to have a better error management.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Notes
* Thanks mostly to StackOverflow and Medium for the endless source of knowledge.
* And thanks for reading up until the end :)
© 2021 GitHub, Inc.
Terms
Privacy
Security
Status
Docs
Contact GitHub
Pricing
API
Training
Blog
About
