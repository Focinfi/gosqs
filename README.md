### SQS

Build a simple queue service.

Plan:

1. A simple queue service run on a single machine, decoupling dependency of storage using interfaces.
2. Build friendly API for common scenario, without guarantee for order and at-most-once push, just for heigh concurrency and simplicity.
3. To be distributed for horizontal scaling, maybe need learn nsq.
4. Support FIFO and at-most-once .
  