# Distributed-Messaging

This repo houses two implementations of a decentralized/distributed console based messaging application.

chat.proto defines the services and messages used by both implementations.

client and client_server implement a messaging service where each client hosts their own server. Clients that want to chat together join the same server and that server handles the communication.

client2 and client_server2 implement a messaging service where each client hosts a server. However, instead of one client hosting server handling a group's communications, each client connects to every server in the communication meaning no single server is responsible for the entire group.

![image](https://user-images.githubusercontent.com/43715044/53779798-733f2980-3eb6-11e9-93fa-454a5e026873.png)

![image](https://user-images.githubusercontent.com/43715044/53779962-1e4fe300-3eb7-11e9-8010-5ec87ea66d45.png)
