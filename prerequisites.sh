#PreRequisites:
#Create Ubuntu VM : 18.04.1-Ubuntu

echo "----------------------------------------------"
echo "SETTING UP PREREQUSITES FOR HYPERLEDGER FABRIC"
echo "----------------------------------------------"

#Get git and install git version 2.17.1
echo "Getting and installing git"
echo "--------------------------"
#cd ~
sudo apt update
sudo apt install git -y
git --version
echo " "
echo " "
sleep 2

#Get go and install go1.13.15
echo "------------------------"
echo "Getting go and installing go1.13.15"
echo "------------------------------------"
#cd ~
wget https://go.dev/dl/go1.13.15.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.13.15.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
source ~/.bashrc
go version
echo " "
rm go1.13.15.linux-amd64.tar.gz

sleep 2

#Get go and install docker-20.10.7
echo "----------------------------"
echo "Getting and installing  docker-20.10.7"
echo "------------------------------------------"
# cd ~
sudo apt update -y
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y
sleep 5
sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
sudo apt update
sudo apt-cache policy docker-ce
sudo apt install docker-compose -y
sleep 5
sudo docker -v

sudo docker-compose -v
echo " "
sleep 5


#Get go and install Node v16.17.0
echo "----------------------------"
echo "Getting go and installing Node v16.17.0"
echo "---------------------------------------"
#cd ~
sudo apt install nodejs -y
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
source ~/.bashrc
nvm install v16.17.0
node -v

sleep 5


#Get go and install Hyperledger binaries 2.2.2 1.4.9
echo "-----------------------------------------------"
echo "Getting go and installing Hyperledger binaries 2.2.2 1.4.9"
echo "------------------------------------------------------------"
#cd ~
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.2 1.4.9

