# analytics-service

Analytics service simple inserts batchs of events in db. 

It has a set of predefines commands it task file for building, testing.


There are branches with different tactics of inserting events and benchmarks using ddosify:

* Wait for ClickHouse response
RESULT
-------------------------------------
Success Count:    5000  (100%)
Failed Count:     0     (0%)

Durations (Avg):
  DNS                  :0.0000s
  Connection           :0.0000s
  Request Write        :0.0000s
  Server Processing    :0.0706s
  Response Read        :0.0000s
  Total                :0.0707s

Status Code (Message) :Count
  200 (OK)    :5000


* Immediately return after body read
RESULT
-------------------------------------
Success Count:    5000  (100%)
Failed Count:     0     (0%)

Durations (Avg):
  DNS                  :0.0000s
  Connection           :0.0000s
  Request Write        :0.0000s
  Server Processing    :0.0012s
  Response Read        :0.0001s
  Total                :0.0013s

Status Code (Message) :Count
  200 (OK)    :5000


* Use pipelining without waiting for ClickHouse response

RESULT
-------------------------------------
Success Count:    5000  (100%)
Failed Count:     0     (0%)

Durations (Avg):
  DNS                  :0.0000s
  Connection           :0.0001s
  Request Write        :0.0000s
  Server Processing    :0.0009s
  Response Read        :0.0001s
  Total                :0.0010s

Status Code (Message) :Count
  200 (OK)    :5000

* Use pipelining with waiting for ClickHouse response
RESULT
-------------------------------------
Success Count:    5000  (100%)
Failed Count:     0     (0%)

Durations (Avg):
  DNS                  :0.0000s
  Connection           :0.0001s
  Request Write        :0.0000s
  Server Processing    :0.0344s
  Response Read        :0.0000s
  Total                :0.0345s

Status Code (Message) :Count
  200 (OK)    :5000

