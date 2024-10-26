# Forum Project - Security

## Description
Forum project unites together 3 mayor parts on webpage development.
1. Server side - build in Golang;
2. Front side - using HTML, CSS and JS commands;
3. Database - SQLITE database used in this project.

Current project manipulates data contained in database in order to create a site and displaying the information.
Project displays information through several data visualizations options and focuses on the creation of events/actions for data visualization.
Project reads only minimal necessary information to present full content of current page. Additional data is loaded from server when Events/actions are executed.
Project consists of user profile creation, information sharing between different users who can access the page and user reactions about different information.

## When you have not started with Forum Project yourself
It is possible that you are missing some vital components in VSCode and get some error messages while trying to run the project.
In my case I had to run these 2 commands to update VSCode
```
sudo apt update
sudo apt install build-essential
```


## How To Run
1. Clone the repo ```git clone https://01.kood.tech/git/anomm/forum```
2. Run the following command: ```go run . ```.
3. Allow access through firewall
4. Open a web browser and navigate to [localhost:8081](https://localhost:8081)
5. Ignore web browser warning about page not being trustworthy and procede to the page. This is because page uses https (Hypertext Transfer Protocol Secure). For more information read [https://www.cloudflare.com/learning/ssl/what-is-https/](https://www.cloudflare.com/learning/ssl/what-is-https/)
6. Before Docker testing please terminate the server.
7. To terminate the server click ```CTRL + "C"```

## Testing
1. Reference for testing can by found in [Audit page](https://01.kood.tech/git/root/public/src/branch/master/subjects/forum/security/audit.md). 

For sql tasks we recommend to look also this [YouTube video](https://www.youtube.com/watch?v=VKg1Dnz7GN0&t=54s). The tables that you are interested are users, topic_header, topic_content.
```
select * from users
select * from topic_header
select * from topic_content
```


For Docker tasks don't forget to open docker software before executing the commands. For generating docker and running it you can also use script:
```
sh docker-script.sh
```
The script will dry to update your docker and asks for password to do so.
The script will for a name that you wish to give to you docker container.

Current database has one user already created
```
UserName: admin@admin.admin
Password: adminadmin
```

2. To test error message capability go to page [localhost:8080/whateverEverYouLike](http://localhost:8080/ThisIsCompletelyDifferentThanYouExpected). 

## Implementation
User creation (Current project allows only Password login)
Login (Password will be crypted. Current project allows only Password login)
Creating Forum topics
Creating Comments under topic
Updating Topics and Comments (only for record creator) 
Deleting Topics and Comments (only for record creator)
Liking and Disliking posts (only when logged in)
Implemented Hypertext Transfer Protocol Secure (HTTPS) protocol 
Implemented of Rate Limiting (You are allowed to have 10 calls in 30 seconds)
Image upload to forum Topics

## Authors
Ainar Nõmm (anomm)
Lauri Laretei (laurilaretei)
Lisann Saaremets (LizAshwood)
Markus Kasemetsa (MarkusKa)
Ragnar Küüsmaa (Ragnar)
