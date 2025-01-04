## Run
- clone this repo
- write .env file (see .envexample file)
- run `sudo docker-compose run -d`

## Usage
you can use any utility for sending json, in examples i will use curl

### Create account
`curl -X POST *your address*/register -d '{"username": "*username*", "password": "*password*"}'`

now we are need to get token

### Login (get token)
`curl -X POST *your address*/login -d '{"username": "*username*", "password": "*password*"}'`

you will get response from server like:

`"token": "*long string*"`

save it.

### Save link
`curl -X POST *your address*/shorten -d '{"alias": "*name for your link*", "link": "*long link*"}' -H "Authorization: Bearer *your token*"`

now, you can write to the browser "*your address*"/*some alias* and you will redirected to the link

### Delete link
`curl -X POST *your address*/*some alias*?action=delete -H "Authorization: Bearer *your token*"`

## Information
tokens will be unusable after 2 hours. to get new token you need to login again

idk how do you will use this shit, but... i think its may be a service for microservice app?
