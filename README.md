



### Example: bind www2.example.com to your public ip address
#### 1. create a file named .env with the following content:
```
CF_DOMAIN=example.com
CF_TOKEN=myyuSxxxxxxxxxxxxxxxxtecgJiGIg

PUB_BIND=www2

```
#### 2. run the executable in linux/unix
./cloudflare-ddns-*


## list all your domains
./cloudflare-ddns-* list
## list all your domains filtered by an addtional string
./cloudflare-ddns-* list www

# Improvements:
ipv6 support


## openwrt certificate issue
### 2021/07/19 09:28:00 ListZonesContext command failed: HTTP request failed: Get "https://api.cloudflare.com/client/v4/zones?name=swiftducks.com&per_page=50": x509: certificate signed by unknown authority
Looks like system on your router lacks CA certificates and can't verify TLS connections on its own. If it's something like OpenWRT, you can install CA bundle like this:
```
opkg update && opkg install ca-bundle
```

Alternatively, you may just put CA cacert.pem from here into one of these locations:
```
/etc/ssl/certs/ca-certificates.crt
/etc/pki/tls/certs/ca-bundle.crt
/etc/ssl/ca-bundle.pem
/etc/pki/tls/cacert.pem
/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem
/etc/ssl/cert.pem
```
