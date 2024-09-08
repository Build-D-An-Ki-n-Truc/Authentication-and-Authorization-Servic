# README

## Introduction
This README provides an overview of the auth API and its endpoints.

## Endpoints

### 1. POST /auth/login/
- Description: Login .
- Data:{
		"username": username,
		"password": password,
	},
- Response: Login status.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/auth/blob/main/image/login.png)

### 2. POST /auth/register/user?crypted=bool
- Description: Register a new user.
- Parameter: crypted=bool (true if password has already been ecrypted with bcrypt, false otherwise)
- Data:{
		"username" : username,
        "password" : password,
        "name": name,
        "email": email, 
        "role":"user",
        "phone": phone,
        "isLocked": false
	},
- Response: Register status.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/auth/blob/main/image/register.png)

### 3. POST /auth/sendOTP
- Description: Send OTP via email to new user to confirm register, This will response with the OTP that has been sent to user.
- Data:{
		"email" : email,
	},
- Response: Email send status.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/auth/blob/main/image/sendEmail.png)

### 4. POST /auth/register/brand?crypted=bool
- Description: Register a new user.
- Parameter: crypted=bool (true if password has already been ecrypted with bcrypt, false otherwise)
- Data:{
		"username" : username,
        "password" : password,
        "name": name,
        "email": email, 
        "role":"user",
        "phone": phone,
        "isLocked": false,
		"brandID": brandID (string)
	},
- Response: Register status.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/auth/blob/main/image/registerBrand.png)
