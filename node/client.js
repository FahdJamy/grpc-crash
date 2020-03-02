'use strict';

const PROTO_PATH = './pb/messages.proto';
const fs = require('fs');
const grpc = require('grpc');
const process = require('process');

const serviceDef = grpc.load(PROTO_PATH);
const PORT = 9002;

const SERVER_ADDRESS = 'localhost:9002';
const EmployeeService = grpc.load(PROTO_PATH).EmployeeService;
const goClient = new EmployeeService(SERVER_ADDRESS, grpc.credentials.createSsl(caCert, key, cert))

// pull in the certificates here
const caCert = fs.readFileSync('certs/ca.crt');
const cert = fs.readFileSync('cert.pem');
const key = fs.readFileSync('key.pem');

const kvpair = {
  'private-key': key,
  'cert_chain': cert,
}

const creds = grpc.credentials.createSsl(caCert, key, cert);
const client = new serviceDef.EmployeeService(`localhost:${PORT}`, creds);

const option = parseInt(process.argv[2], 10);
switch (option) {
  case 1:
    sendMetadata(client);
    break;
  case 2:
    getByBadgeNumber(client);
    break;
  case 3:
    getAll(client);
    break;
  case 4:
    addPhoto(client);
    break;
  case 5:
    saveAll(client);
    break;
}

function sendMetadata(client) {
  const md = new grpc.Metadata();
  md.add('username', 'googleUser');
  md.add('password', 'krs1krs1');

  client.getByBadgeNumber({}, md, function() {});
}

// unary operarion
function getByBadgeNumber(client) {
  client.getByBadgeNumber({badgeNumber: 2080}, function (err, res) {
    if (err) {
      console.log(err);
    } else {
      console.log(res.employee);
    }
  })
}

// getAll, send a request and receive a stream of responses -------> serverSideStreaming
function getAll(client) {
  const call = client.getAll({});

  call.on('data', function(data) {
    console.log(data.employee);
  })
}

// send a stream of requests and recieve one reponse -------> clientSideStreaming
function addPhoto(client) {
  const md = new grpc.Metadata();
  md.add('badgenumber', '2080');
  const call = client.addPhoto(md, function (err, res) {
    if (err) {
      console.log(err);
    } else {
      console.log(res.employee);
    }
  });

  const stream = fs.createReadStream('photo.jpg');
  stream.on('data', function (chunk) {
    call.write({data: chunk});
  });
  stream.on('end', function () {
    call.end();
  })
}

// send a stream of requests and recieve a stream of responses
function saveAll(client) {
  const employees = [
    {
      id: 4,
      badgeNumber: 2030,
      firstName: "dont",
      lastName: "go",
      vacationAccrualRate: 1,
      vacationAccrued: 3,
    },
    {
      id: 5,
      badgeNumber: 2040,
      firstName: "home",
      lastName: "groot",
      vacationAccrualRate: 1,
      vacationAccrued: 40,
    }
  ];

  const call = client.saveAll();
  call.on('data', function (emp) {
    console.log(emp);
  });
  employees.forEach(function(emp) {
    call.write({employee: emp});
  });
  call.end();
} 
