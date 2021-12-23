
const process = require ("process");
const http = require ("http");


const _dummyBodyString = "hello world!\n";
const _dummyBodySize = _dummyBodyString.length;
const _dummyBodyBuffer = Buffer.alloc (_dummyBodySize, _dummyBodyString, "utf-8");
const _dummyHeaders = [
		"Content-Length", _dummyBodySize.toString (),
		"Content-Type", "text/plain; charset=utf-8",
		"Content-Encoding", "identity",
		"Cache-Control", "no-store, max-age=0",
	];


function _handler (_request, _response) {
		_response.sendDate = false;
		_response.writeHead (200, "OK", _dummyHeaders);
		_response.end (_dummyBodyBuffer);
	};

const _server = http.createServer ({}, _handler);


var _endpoint_ip = "127.0.0.1";
var _endpoint_port = 8080;
switch (process.argv.length) {
	case 2 :
		break;
	case 3 :
		_endpoint_port = parseInt (process.argv[2]);
		break;
	case 4 :
		_endpoint_ip = process.argv[2];
		_endpoint_port = parseInt (process.argv[3]);
		break;
	default :
		throw Error ("[a102837d]");
}

_server.listen (_endpoint_port, _endpoint_ip, 65536);

