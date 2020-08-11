# Secretly-cli
This is a CLI tool wrapped around [my Secretly package](https://github.com/GoZaddy/secret.ly)

## Storing a secret
secret.ly-cli set --name "test_secret" --value "test_secret_value" --key "encryption_key"

## Retrieving a secret
secret.ly-cli get --name "test_secret" --key "encryption_key"
