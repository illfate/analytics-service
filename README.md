# analytics-service

Analytics service simple inserts batchs of events in db. 

It has a set of predefines commands it task file for building, testing.


There are branches with different tactics of inserting events.

* Wait for ClickHouse response
* Immediately return after body read
* Use pipelining without waiting for ClickHouse response
* Use pipelining with waiting for ClickHouse response
