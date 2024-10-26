sudo apt-get update
sudo apt-get install dos2unix
dos2unix docker-script.sh 
sudo docker build -t forum-docker .
echo ""
echo "Type a name for your forum container: " 
read name
echo ""
sudo docker run -d -p 8081:8081 --name $name forum-docker
echo ""
echo "Server was started on PORT :8081. Please go to your web browser and open https://localhost:8081"
echo ""
echo "To stop the docker run: "
echo "$ sudo docker stop $name"
echo ""
echo "To delete docker first stop it, then run: "
echo "$ sudo docker rm $name"