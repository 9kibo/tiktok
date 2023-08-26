# gorm

## trace error

```json
[
  {
    "data": {
      "elapsed": 3.6784,
      "err": {
        "Number": 1241,
        "SQLState": [
          50,
          49,
          48,
          48,
          48
        ],
        "Message": "Operand should contain 1 column(s)"
      },
      "file": "A:/code/backend/go/tiktok/biz/dao/message_page_no_test.go:216",
      "rows": 0,
      "sql": "SELECT * FROM `message` WHERE from_user_id in (1) and to_user_id = (5410,4551,1821,1051,4937,3320,5758,2148,3216,5449,4084,5287,2574,4836,1515,3873,5968,5091,1790,4331,4273,5956,5911,1637,2378,2109,3072,2650,2588,4971,4344,5443,2623,4459,2803,2739,2883,4907,4932,3780,4115,3716,1020,4071,5652,2907,3719,5062,4110) and id = (SELECT id FROM message as m2 WHERE from_user_id = message.from_user_id  and to_user_id = 1 ORDER BY created_at desc LIMIT 1)"
    },
    "level": "debug",
    "module": "GORM",
    "msg": "trace error",
    "time": "2023-08-24 22:13:12"
  },
  {
    "data": {
      "elapsed": 14.1054,
      "err": {},
      "file": "A:/code/backend/go/tiktok/biz/dao/message_page_no_test.go:203",
      "rows": 0,
      "sql": "SELECT * FROM `message` WHERE from_user_id=1 and to_user_id in (5410,4551,1821,1051,4937,3320,5758,2148,3216,5449,4084,5287,2574,4836,1515,3873,5968,5091,1790,4331,4273,5956,5911,1637,2378,2109,3072,2650,2588,4971,4344,5443,2623,4459,2803,2739,2883,4907,4932,3780,4115,3716,1020,4071,5652,2907,3719,5062,4110) and id = (SELECT id FROM message as m2 WHERE from_user_id = 1 and to_user_id = message.to_user_id ORDER BY created_at desc LIMIT 1)"
    },
    "level": "debug",
    "module": "GORM",
    "msg": "trace error",
    "time": "2023-08-24 22:13:12"
  }
]
```

## trace warn

```json
[]
```

## trace info

```json
[
  {
    "data": {
      "elapsed": 2.1677,
      "file": "A:/code/backend/go/tiktok/biz/dao/message_page_no_test.go:240",
      "rows": 1,
      "sql": "SELECT * FROM `message` WHERE (from_user_id = 1 and to_user_id =4937) or (to_user_id = 4937 and from_user_id =1) ORDER BY created_at desc LIMIT 1"
    },
    "level": "debug",
    "module": "GORM",
    "msg": "trace info",
    "time": "2023-08-24 22:13:12"
  }
]
```