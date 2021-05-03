# grpc - a learning example of CalculatorService

## this will include following 4 message communication scenario 
### unary
Sum API - to implement a Sum RPC Unary API in a CalculatorService
- The function takes a Request message that has two integers, and returns a Response that represents the sum of them.
### Server streaming
PrimeNumberDecomposition API - to implement a PrimeNumberDecomposition RPC Server Streaming API in a CalculatorService
- The function takes a Request message that has one integer, and returns a stream of Responses that represent the prime number decomposition of that number
### Client streaming
ComputeAverage API 
- The function takes a stream of Request message that has one integer, and returns a Response with a double that represents the computed average
### Bi-Direction streaming 
FindMaximum API 
- The function takes a stream of Request message that has one integer, and returns a stream of Responses that represent the current maximum between all these integers

### Error, Deadline

### SSL Encryption 

### CRUD API with Mongodb

- package : calc 
 simple calculator service 
**proto : calc.proto

- package : blog
 simple blog api to create, read, update, delete and listing using streaming. 
 **proto : blog.proto
