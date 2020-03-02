'use strict'

const path = require('path')
const PROTO_PATH = path.join('pb', 'messages.proto')
const SERVER_ADDRESS = 'localhost:50000';

const grpc = require('grpc')
const HelloService = grpc.load(PROTO_PATH).HelloService;

const client = new HelloService(SERVER_ADDRESS, grpc.credentials.createInsecure());

function main () {
  client.sayHello({Name: 'Google'}, function (err, reponse) {
    if(err) {
      console.log(err);
      return
    }
    console.log(reponse.Message)
  })
}

main();