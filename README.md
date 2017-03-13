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
3. [x] Use golang channel to reduce loop times of Pop.

#### Message Storage

1. [ ] Search distributed reliable K/V db.

1. [ ] Time-based Group or ID Generator?
  
  1. Design new message group rule, maybe 100items/group?

  1. Add message id into group must hold a RWLock, try etcd relative api(Compare And Swap?).

  1. Maybe more dependency on databse.

1. [ ] Batch messages addation?

1. [ ] Compact and put the cold data(Message and Group) into disk?

#### ID generator

1. [ ] Maybe use etcd?


#### Client SDK

1. [ ] Replace HTTP with TCP?

2. [ ] Inergrate `ID generator` into SDK.

#### Testing

1. [ ] Use interface mock.

2. [ ] Build benthmark.

#### Load Balancing
1. [ ] Build center/master to collect the all nodes/service status

  1. Node <=> Database
  2. Consumer <=> Node

2. [ ] Shard the key in Database?

#### Admin page
1. [ ] System health monitoring.

1. [ ] User management, `acess_key`, `api_key`.

1. [ ] Agent auth.

1. [ ] Flow control.
