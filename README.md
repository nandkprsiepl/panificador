>PreRequisites:

>Create Ubuntu VM : 18.04.1-Ubuntu

Installation Steps
------------------

git clone https://github.com/MDxBlocksRepo/syntheses-hyperledger

sudo su

cd syntheses-hyperledger

chmod +x *.sh

chmod +x artifacts/channel/*.sh


./prerequisites.sh

Add the following paths to /root/.bashrc by selecting actual path of fabric-samples

export PATH=$PATH:/usr/local/go/bin

export PATH=$PATH:~USER/syntheses-hyperledger/fabric-samples/bin

Start Duplication session

sudo su

nvm install v16.17.0

>node --version
 
v16.17.0


>go version
 
go version go1.13.15 linux/amd64
 
 
>docker version

Client: Docker Engine - Community
 Version:           20.10.17

>docker-compose version

docker-compose version 1.17.1, build unknown

cd syntheses-hyperledger

cd artifacts/channel
 
./create-artifacts.sh

cd ../..

./start.sh
