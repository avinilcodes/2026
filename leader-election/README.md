# CAP Theorem in SpamDetection

SpamDetection is designed as a AP system.

## Consistency
It is not consistent that is fine, meaning that if user gets false positive that is fine.

## Availability
Service must be available all the times, so that users get their responses.

## Partition Tolerance
Network partitions are assumed. All decisions prioritize availability
over consistency during partitions.

## Tradeoff
User gets response every time, even if the consistency is not there.
