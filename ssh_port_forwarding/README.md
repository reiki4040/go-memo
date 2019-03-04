# SSH port forwarding

example SSH port forwarding in Go.

**Many code changed about connection handling. see https://github.com/reiki4040/mogura**

## do it

bastion and http server host starting...

```
vagrant up
  # take some minutes...
```

check your vagrant ssh private key.

```
vagrant ssh-config bastion
  # get your ssh private key path (vagrante generated)
  # if you set http host private key, then panic with nil.
```

replace config.yml

```
name: test vagrant
ssh_bastion_host_port: "localhost:2222"
ssh_user: vagrant
ssh_key_file_path: <YOUR ENV PRIVATE KEY FILE PATH> 
local_bind_port: "localhost:8080"
forwarding_remote_port: "192.168.50.6:80"
```

start port forwarding.

```
./ssh_port_forwarding
```

curl internal http server with port forwarding.

```
curl http://localhost:8080/
  # success if got nginx welcome page html.
```

## refs

- https://orebibou.com/2019/02/golang%E3%81%A7ssh-port-forwarding%E3%82%92%E3%81%99%E3%82%8B/
- https://godoc.org/golang.org/x/crypto/ssh
- https://stackoverflow.com/questions/45441735/ssh-handshake-complains-about-missing-host-key
