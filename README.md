## Run
- clone this repo
- create .env file (check out .envexample file)
- run `sudo docker-compose run -d`

## Usage
you can use any utility for sending jsons, in the examples i use curl

### Create account
`curl -X POST *your address*/register -d '{"username": "*username*", "password": "*password*"}'`

next we are need to get token

### Login (get token)
`curl -X POST *your address*/login -d '{"username": "*username*", "password": "*password*"}'`

you will get a response from the server like:

`"token": "*long string*"`

save it.

### Save link
`curl -X POST *your address*/shorten -d '{"alias": "*name for your link*", "link": "*long link*"}' -H "Authorization: Bearer *your token*"`

now, you can type in the browser "*your address*"/*some alias* and you will redirect to the link

### Delete link
`curl -X POST *your address*/*some alias*?action=delete -H "Authorization: Bearer *your token*"`

## Information
tokens will become unusable after 2 hours. to get a new token you should to login again
