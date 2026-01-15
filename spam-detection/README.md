# CAP Theorem in SpamDetection

SpamDetection is designed as a AP system.

## Consistency
Authentication must be correct. Returning stale or incorrect auth data
is worse than rejecting a request.

## Availability
If the database is unavailable, SpamDetection will reject login/signup
requests instead of serving stale data.

## Partition Tolerance
Network partitions are assumed. All decisions prioritize consistency
over availability during partitions.

## Tradeoff
During DB outages, users may be unable to authenticate,
but no invalid tokens will be issued.
