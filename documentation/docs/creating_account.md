# Creating an account

The following APIs correspond to the following agreed upon workflows;

- https://drive.google.com/file/d/1zRjZse_dyyjVhqVw48BnLRt7U1ovgd0D/view?usp=sharing

- https://drive.google.com/file/d/1LnvGj_jOFMgpi0HNRw8A_zspjH0-SamQ/view?usp=sharing

## APIs

### Check if phone number exists

This first step to creating an account, is to creating whether a the provided phone number has been registered to another user or not.

That is what this API does. If the phone number is unique (unknown), it will send an OTP to the provided phone number

It check both PRIMARY PHONE NUMBERS and SECONDARY PHONE NUMBERS

Example (phone number exists);

```sh
http https://profile-testing.healthcloud.co.ke/verify_phone phoneNumber="+254718376163"


HTTP/1.1 400 Bad Request
Cache-Control: private
Content-Encoding: gzip
Content-Type: text/html
Date: Tue, 05 Jan 2021 12:50:17 GMT
Server: Google Frontend
Transfer-Encoding: chunked
X-Cloud-Trace-Context: e0f23bcbd33ad69e1a87f66b286cd91e;o=1
vary: Accept-Encoding

{
    "error": "4"
}
```

Example (phone number does not exist)

```sh
http https://profile-testing.healthcloud.co.ke/verify_phone phoneNumber="+254715825862"

HTTP/1.1 200 OK
Cache-Control: private
Content-Encoding: gzip
Content-Type: text/html
Date: Tue, 05 Jan 2021 12:55:45 GMT
Server: Google Frontend
Transfer-Encoding: chunked
X-Cloud-Trace-Context: eef2c7afb07866afbcebaff0342c0b51
vary: Accept-Encoding

{
    "otp": "{\"otp\":\"779124\"}"
}
```

!!! note
Go to `Error codes` page for details regarding error codes.
