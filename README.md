# QUIC-VPN

## Note:

increase send and receive buffer

Temporaly

sudo sysctl -w net.core.rmem_max=8388608
sudo sysctl -w net.core.wmem_max=8388608

nano /etc/sysctl.conf

```
net.core.wmem_max=8388608
net.core.rmem_max=8388608
```

sudo sysctl -p