# ACMEWrapper

Add Let's Encrypt support to your golang server in 5 lines of code.

## Testing

Running the tests is a bit of a chore, since it requires a valid domain name, and access to port 443.
This is because ACMEWrapper uses the Let's Encrypt staging server to make sure the code is working.

To test on your own server, you need to change the domain name to your domain, and set a custom testing port
that will be routed to 443:

```bash
sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 8443
export TLSADDRESS=":8443"
export DOMAIN_NAME="example.com"
go test
```

To delete the port forwarding rule, find it in your tables
```bash
sudo iptables -t nat --line-numbers -n -L
```

And delete the number that it was at
```bash
iptables -t nat -D PREROUTING 2
```

(This is based on http://serverfault.com/questions/112795/how-can-i-run-a-server-on-linux-on-port-80-as-a-normal-user)
