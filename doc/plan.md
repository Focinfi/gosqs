## Plan

### Consumer Lopp Rule

Add new-message subscribtion for the consumer which has been waiting for long time.
  1. [ ]Design the value of the "long time".
  2. [ ]Using etcd for subscribtion.

### Message Storage

1. [ ]Search distributed reliable K/V db.

1. [ ]Design new message group rule, maybe 100items/group?

1. [ ]Add message id into group must hold a RWLock, try etcd relative api.

1. [ ]Batch messages addation?

### ID generator

1. [ ]Maybe use etcd?


### Client SDK

1. [ ]Replace HTTP with TCP? 
2. [ ]Inergrate `ID generator` into SDK.

### Testing

1. [ ]Use interface mock.
2. [ ]Build benthmark.
