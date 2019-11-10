Run Simulation in HPS Cluster
---

Solon is the City University high-performance computation cluster.
It is behind the firewall so it should be accessed through university
remote access services: https://vpn.city.ac.uk/workplace/access/home

To send small text files from cluster use mail:

    mail -a /opt/emailfile.eml -s "Email File" alex.netkachov@gmail.com < /dev/null

To send larger files, use files@alexatnet.com account.

    scp ... files@alexatnet.com:/home/files/...

The account "files" is just a basic account created as:

    adduser --gecos "" files
    usermod -a -G ssh-clients files
    mkdir -m 700 ~/.ssh
    vi ~/.ssh/authorized_keys
    ... paste public keys from macbook and solon account
    chmod 600 ~/.ssh/*

On the solon:

    ssh-keygen -t rsa
    vi ~/.ssh/config
    Host 109.74.205.143
      IdentityFile ~/.ssh/id_rsa
      IdentitiesOnly yes
      Port 22

where 109.74.205.143 is the ip address of alexatnet.com

To SSH to solon use the following trick:

    acmp148@solon$ ssh -R 9091:localhost:22 files@alexatnet.com
    files@alexatnet.com$ ssh -p 9091 acmp148@localhost

If you want to establish permanent connection:

    acmp148@solon$ screen
    acmp148@solon$ eval "$(ssh-agent -s)"
    acmp148@solon$ ssh-add
    acmp148@solon$ while [ 1 == 1 ]; do ssh -R 9091:localhost:22 files@alexatnet.com; sleep 10; done
    Ctrl+A Ctrl+D

To copy files from files@alexatnet.com to solon use

    files@alexatnet.com$ scp -P 9091 ... acmp148@localhost:/home/acmp148/...

To view the queue on the cluster use

    showq