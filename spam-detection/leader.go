// Leader is responsible for:
// - rotating JWT keys
// - cleanup expired sessions
// - issuing refresh tokens

/*
How leader is chosen
Leader is chosen via a distributed lock.

What leader does
The leader is responsible for:
- rotating JWT keys
- cleaning up expired sessions
- issuing refresh tokens


How leader maintains its leadership
The leader maintains its leadership by periodically renewing the distributed lock.
If it fails to renew the lock, another instance can take over as leader.
How leader handles failover
If the current leader fails to renew the lock, another instance can acquire the lock and become the new leader.
*/