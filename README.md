### SQS

Build a simple queue service.

### Design

1. Simple
1. Distributed
2. Most-at-once Delivery
3. Ordered Message
4. Horizontal Scaling

### Status

Basic implementation.

### Plan

#### Consumer Lopp Rule

Add new-message subscribtion for the consumer which has been waiting for long time.

1. [x] Design the value of the "long time".
2. [x] Use etcd watcher for subscribtion.
3. [ ] Use redis watcher to reduce loop times of Pop.

#### Message Storage

1. [ ] Search distributed reliable K/V db.

1. [ ] Design new message group rule, maybe 100items/group?

1. [ ] Add message id into group must hold a RWLock, try etcd relative api.

1. [ ] Batch messages addation?

#### ID generator

1. [ ] Maybe use etcd?


#### Client SDK

1. [ ] Replace HTTP with TCP? 
2. [ ] Inergrate `ID generator` into SDK.

#### Testing

1. [ ] Use interface mock.
2. [ ] Build benthmark.
