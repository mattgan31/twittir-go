# Twittir Golang API

## About this project
This is an API for a Social Media application. I built this API with the aim of both recalling and honing my Golang skills.

## How to install ?

### 1. You must clone to your local repository
To clone this repository to your local machine, make sure you have git installed, then run the following command:
```sh
git clone https://github.com/mattgan31/twittir-go.git
```

### 2. Settings the database file
Next, you can configure the database settings in the database/database.go file. Adjust them according to your desired database.

### 3.  Next you can start the app
To run the application, you need to install golang first. After that, you can run the following command in the terminal to start the application:
```sh
go run main.go
```

## Endpoint Documentation

## Auth
### - Register User
#### Request
```sh
# POST | localhost:3001/api/register
{
    "email": "your_email@email.com",
    "username": "your_username",
    "fullname": "Your Fullname",
    "password": "yourpassword"
}
```
#### Response
```json
{
    "data":
    {
        "email": "your_email@email.com",
        "fullname": "Your Fullname",
        "id": 3
    },
    "status": "success"
}
```

### - Login User
#### Request
```sh
# POST | localhost:3001/api/login
{
    "username": "your_username",
    "password": "yourpassword",
}
```

#### Response
```json
{
    "token": "GENERATEDTOKEN"
}
```

## Create Post & Comment
### - Create Post
#### Request
```sh
# POST | localhost:3001/api/posts
{
    "post": "Your Post Description"
}
```

#### Response
```json
{
    "post":
    {
        "createdAt": "DATETIMEZ",
        "id": 1,
        "post": "Your Post Description",
        "user_id": 1
    }
}
```

### - Create Comment
#### Request
```sh
# POST | localhost:3001/api/posts/[postID]/comment
{
    "description": "Your Comment Description"
}
```

#### Response
```json
{
    "post":
    {
        "createdAt": "DATETIMEZ",
        "id": 1,
        "post": "Your Comment Description",
        "user_id": 1
    }
}
```

## Likes
### - Create Like Post
#### Request
```sh
# POST | localhost:3001/api/posts/[postID]/like
```

#### Response
```json
{
    "message": "Like post with id [postID] success"
}
```

### - Create Like Comment
#### Request
```sh
# POST | localhost:3001/api/comments/[commentID]/like
```

#### Response
```json
{
    "message": "Like comment with id [commentID] success"
}
```

## Get Data Posts
### - Get All Posts
#### Request
```sh
# GET | localhost:3001/api/posts
```

#### Response
```json
{
    "posts":
    [
        {
            "id": 1,
            "post": "Your Description Post",
            "created_at": "DATETIMEZ",
            "user":{
                "id": 1,
                "username":"your_username"
            },
            "likes": [],
            "comments": [],
        }
    ]
}
```

### - Show Posts
#### Request
```sh
# GET | localhost:3001/api/posts/[postID]
```

#### Response
```json
{
    "posts":
    {
        "id": 1,
        "post": "Your Description Post",
        "created_at": "DATETIMEZ",
        "user":{
            "id": 1,
            "username":"your_username"
        },
        "likes": [],
        "comments": [],
    }
}
```

### - Get All Posts By User
#### Request
```sh
# GET | localhost:3001/api/posts/user/[userID]
```

#### Response
```json
{
    "posts":
    [
        {
            "id": 1,
            "post": "Description Post",
            "created_at": "DATETIMEZ",
            "user":{
                "id": 1,
                "username":"username"
            },
            "likes": [],
            "comments": [],
        }
    ]
}
```

## Delete Post/Comment
### Delete Post
#### Request
```sh
# DELETE | localhost:3001/api/posts/[postID]
```

#### Response
```json
{
    "message": "Post with id [postID] deleted successfully"
}
```

### Delete Comment
#### Request
```sh
# DELETE | localhost:3001/api/comments/[commentID]
```

#### Response
```json
{
    "message": "Comment with id [commentID] deleted successfully"
}
```

## Show Users
### - Show Profile
#### Request
```sh
# POST | localhost:3001/api/users/profile
```

#### Response
```json
{
    "fullname": "Your Fullname",
    "id": 1,
    "ProfilePicture": "",
    "username": "your_username"
}
```

### - Show Other Profile User
#### Request
```sh
# POST | localhost:3001/api/users/[userID]
```

#### Response
```json
{
    "fullname": "Fullname",
    "id": 1,
    "ProfilePicture": "",
    "username": "username"
}
```

## Thank you
I will continue to update this application as it still lacks some features compared to existing social media applications.
