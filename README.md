# typebot-telegram-adapter

## Deploying

(requires make, podman-docker and yarn):

```
# aws ecr get-login-password --region us-west-2 | podman login -u AWS --password-stdin "235368163414.dkr.ecr.us-west-2.amazonaws.com"
$ go get -v ./...
$ yarn install --frozen-lockfile
$ make && yarn sls deploy
```

## Details

