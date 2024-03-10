API v3
===============

## Links
[How to handle multiple collections in neo4j](https://stackoverflow.com/questions/28099125/how-to-unwind-multiple-collections)
[Grpc buf](https://connectrpc.com/docs/web/using-clients)
[How to setup grpc-webproxy](https://github.com/grpc-ecosystem/grpc-gateway)

## Docker on AWS EC2
```bash
To get Docker running on the AWS AMI you should follow the steps below (these are all assuming you have ssh'd on to the EC2 instance).

Update the packages on your instance

[ec2-user ~]$ sudo yum update -y

Install Docker

[ec2-user ~]$ sudo yum install docker -y

Start the Docker Service

[ec2-user ~]$ sudo service docker start

Add the ec2-user to the docker group so you can execute Docker commands without using sudo.

[ec2-user ~]$ sudo usermod -a -G docker ec2-user

You should then be able to run all of the docker commands without requiring sudo. After running the 4th command I did need to logout and log back in for the change to take effect.
```
