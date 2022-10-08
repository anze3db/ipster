# ipster

`ipster` is a command line tool that keeps your CloudFlare DNS record in sync with your Raspberry Pi's IP 🤝

## Install and Run

```
wget https://github.com/anze3db/ipster/releases/download/22.1/ipster.pi
chmod +x ./ipster.pi
./ipster.pi
```

`ipster` will need access to your CloudFlare account so you'll have to create an [API_TOKEN](https://dash.cloudflare.com/profile/api-tokens) (use the Edit Zone DNS template).

Run with the necessary environment variables:
```
IPSTER_CLOUDFLARE_API_TOKEN=xxxxxxxxx_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy IPSTER_CLOUDFLARE_ZONE_NAME=example.com IPSTER_CLOUDFLARE_DNS_RECORD_NAME=home.example.com ./ipster.pi
```

Required environmental variables:

* `IPSTER_CLOUDFLARE_API_TOKEN` - your CloudFlare [API_TOKEN](https://dash.cloudflare.com/profile/api-tokens) (Use the Edit zone DNS template).
* `IPSTER_CLOUDFLARE_ZONE_NAME` - your CloudFlare zone name. Usually your domain name e.g. `example.com`.
* `IPSTER_CLOUDFLARE_DNS_RECORD_NAME` - the CloudFlare DNS record that you want to keep in sync e.g. `home.example.com`.

## Example run

`ipster` will compare your machine's IP address to the IP address in the DNS Record in CloudFlare and update the DNS Record if the two don't match.

```
IPSTER_CLOUDFLARE_API_TOKEN=xxxxxxxxx_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy IPSTER_CLOUDFLARE_ZONE_NAME=example.com IPSTER_CLOUDFLARE_DNS_RECORD_NAME=home.example.com ipster.pi
2022/10/08 23:17:11 Verifying IPs
2022/10/08 23:17:13 No change
2022/10/08 23:18:11 Verifying IPs
2022/10/08 23:18:13 No change
2022/10/08 23:19:11 Verifying IPs
2022/10/08 23:19:13 IPs do not match. Updating...
2022/10/08 23:19:14 DNS Record updated!
2022/10/08 23:20:11 Verifying IPs
2022/10/08 23:20:13 No change
```

# Compile for Raspberry Pi

```
env GOOS=linux GOARCH=arm GOARM=5 go build
```
