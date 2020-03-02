'use strict';

const PROTO_PATH = './pb/messages.proto';
const fs = require('fs');
const grpc = require('grpc');

const serviceDef = grpc.load(PROTO_PATH);
const PORT = 9000;

// pull in the certificates here
const caCert = fs.readFileSync('certs/ca.crt');
const cert = fs.readFileSync('certs/server.crt');
const key = fs.readFileSync('certs/server.key');
const employees = require('./employee').employees;
const kvpair = {
  'private-key': key,
  'cert_chain': cert,
}
const creds = grpc.ServerCredentials.createSsl(caCert, [kvpair]);
const server = new grpc.Server();

server.addProtoService(serviceDef.EmployeeService.service, {
  getByBadgeNumber: getByBadgeNumber,
  addPhoto: addPhoto,
  saveAll: saveAll,
  getAll: getAll,
  save: save,
});
server.bind(`0.0.0.0:${PORT}`, creds);
console.log('Running on port ', PORT);
server.start()

// unary
function getByBadgeNumber(call, callback) {
  const badgeNumber = call.request.badgeNumber;
  for (let i = 0; i < employees.length; i++){
    if (employees[i].badgeNumber === badgeNumber) {
      callback(null, {employee: employees[i]});
      return;
    }
  }
  callback('error')
}

// client streaming
function addPhoto(call, callback) {
  const md = call.metadata.getMap();
  for (let key in md) {}

  let result = new Buffer(0);
  call.on('data', function(data) {
    result = Buffer.concat([result, data.data]);
  });
  call.on('end', function() {
    callback(null, {isOk: true});
  });
}

function save(call, callback) {}

// Bidirectional streaming
function saveAll(call) {
  call.on('data', function(emp) {
    employees.push(emp.employee);
    call.write({employee: emp.employee});
  });
  call.on('end', function() {
    call.end();
  });
}

// server streaming in node JS
function getAll(call) {
  employees.forEach(function(emp) {
    call.write({employee: emp});
  });
  call.end()
}
