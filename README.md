# Distributed-Messaging

This repo houses two implementations of a decentralized/distributed console based messaging application.

chat.proto defines the services and messages used by both implementations.

client and client_server implement a messaging service where each client hosts their own server. Clients that want to chat together join the same server and that server handles the communication.

client2 and client_server2 implement a messaging service where each client hosts a server. However, instead of one client hosting server handling a group's communications, each client connects to every server in the communication meaning no single server is responsible for the entire group.

