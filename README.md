



### Example: bind www2.example.com to your public ip address
#### 1. create a file named .env with the following content:
```
CF_DOMAIN=example.com
CF_TOKEN=myyuSxxxxxxxxxxxxxxxxtecgJiGIg

PUB_BIND=www2

```
#### 2. run the executable in linux/unix
./cloudflare-ddns


## list all your domains
./cloudflare-ddns list
## list all your domains filtered by an addtional string
./cloudflare-ddns list www

# Improvements:
ipv6 support