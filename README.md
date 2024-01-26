# github-webhook

The program consists of a github webhook event listener that executes custom commands based on user configuration.

## Getting started

First, you will need to download the binary program and make it executable:

```sh
# Downloading with curl
sudo curl -L -o /usr/bin/github-webhook https://github.com/zanz1n/github-webhook/releases/latest/download/github-webhook

# Using chmod to change file permissions
sudo chmod u+x /usr/bin/github-webhook
```

To run the program:

```sh
sudo github-webhook --config /path/to/the/config.yml
```

## Configuration

An example of the configuration file (with comments) can be found [here](https://github.com/zanz1n/github-webhook/blob/main/config.example.yml).
