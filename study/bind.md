##### bind

激活设备,APP绑定设备连云端,设备必须在线
授权,开发者注册设备,不一定已授权,也不一定练过云端
设备导入但是不一定连接云端




```
mysql> select * from 5473_device_info where sub_domain != "5737" limit 100;
+------+------------+--------------+--------------------+------------------+--------+---------+---------------------+-------------+
| did  | sub_domain | device_id    | name               | aes_key          | status | reserve | create_time         | modify_time |
+------+------------+--------------+--------------------+------------------+--------+---------+---------------------+-------------+
|  318 | 5864       | F0FE6B792CC8 | JST-R506净水机     | rGxjjqPm1ph0J0G5 |      1 | NULL    | 2017-09-06 03:26:35 | NULL        |
|  631 | 5864       | F0FE6B792305 |                    | witsQ84yM9e9nQMP |      1 | NULL    | 2017-09-05 11:00:07 | NULL        |
|  661 | 5864       | F0FE6B791F5F | JST-R506净水器     | ZR6eTbH44zRUO3eS |      1 | NULL    | 2017-09-05 11:50:34 | NULL        |
|  742 | 5864       | F0FE6B919A38 | R505               | nUOUEZ0RM8oKW6Na |      1 | NULL    | 2017-10-22 09:11:23 | NULL        |
| 2023 | 5864       | F0FE6B919A21 |                    | TLfCzGXpMc17Lyx3 |      1 | NULL    | 2017-09-25 05:18:36 | NULL        |
| 2024 | 5864       | F0FE6B91998B | JST-R505净水器     | DODo3vNZXeuAUGBL |      1 | NULL    | 2017-09-20 01:30:48 | NULL        |
| 2132 | 5864       | F0FE6B919B24 | JST-R505净水机     | U85ZFPn6Z4crwOk2 |      1 | NULL    | 2017-10-13 10:08:54 | NULL        |
| 2179 | 5864       | F0FE6B919994 | 碧云泉             | krOaRifpyEWHKwqk |      1 | NULL    | 2017-10-21 10:30:59 | NULL        |
| 2227 | 5738       | F0FE6B792DFC |                    | pLN9U7rvc4g2oZ4z |      1 | NULL    | 2017-09-23 07:14:55 | NULL        |
| 2228 | 5738       | F0FE6B91990A | JST-R702净水器     | xE7qXEdieRHHnGku |      1 | NULL    | 2017-09-23 07:10:09 | NULL        |
| 2251 | 5864       | F0FE6B919ADD | JST-R505净水器     | ZM0CynVOErAs6kU0 |      1 | NULL    | 2017-09-30 11:04:22 | NULL        |
| 2254 | 5864       | F0FE6B919853 | 健康常在           | ZYIFdh2KHGpzZN8F |      1 | NULL    | 2017-10-23 02:08:20 | NULL        |
| 2258 | 5864       | F0FE6B9194AF | JST-R505净水机     | 0uKuOW6W9TqmE1NS |      1 | NULL    | 2017-10-23 05:08:50 | NULL        |
| 2264 | 5864       | F0FE6B9199B6 | 小云云             | OX0J1U8bLSRKQpVf |      1 | NULL    | 2017-10-22 03:44:10 | NULL        |
| 2272 | 5864       | F0FE6B9196DB | 可可的小管家       | 8o5orJ0ghYaYp8Dw |      1 | NULL    | 2017-10-20 11:34:36 | NULL        |
| 2300 | 5864       | F0FE6B91958D | JST-R505净水机     | IZqS8gxb3L86qbgX |      1 | NULL    | 2017-10-23 03:44:52 | NULL        |
| 2392 | 5864       | F0FE6B919A43 | JST-R505净水机     | QSGKQQXDHdi2uRZ2 |      1 | NULL    | 2017-10-01 14:36:59 | NULL        |
| 2715 | 5864       | F0FE6B9195B2 | 玫瑰金             | z5EDkc8D2EUgqG6E |      1 | NULL    | 2017-10-20 10:55:29 | NULL        |
+------+------------+--------------+--------------------+------------------+--------+---------+---------------------+-------------+
 
 
绑定失败,子域不匹配
sql书写



curl -X POST \
  http://10.0.0.6:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:kingclean" \
  -H "x-zc-major-domain-id:5473" \
  -H "x-zc-sub-domain:r506" \
  -H "x-zc-sub-domain-id:5737" \
  -d "{\"columnName\":\"deviceType\",\"columnType\":4,\"columnLength\":10}"



curl -X POST \
  http://10.0.0.6:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:demo_ac" \
  -H "x-zc-major-domain-id:995" \
  -H "x-zc-sub-domain:ableair" \
  -H "x-zc-sub-domain-id:1003" \
  -d "{\"columnName\":\"deviceType\",\"columnType\":4,\"columnLength\":10}"


curl -X POST \
  http://10.51.67.4:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:kingclean" \
  -H "x-zc-major-domain-id:5473" \
  -H "x-zc-sub-domain:r506" \
  -H "x-zc-sub-domain-id:5737" \
  -d "{\"columnName\":\"deviceType\",\"columnType\":4,\"columnLength\":10}"


curl -X POST \
  http://10.51.67.4:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:kingclean" \
  -H "x-zc-major-domain-id:5473" \
  -H "x-zc-sub-domain:r506" \
  -H "x-zc-sub-domain-id:5737" \
  -d "{\"columnName\":\"deviceSnCode\",\"columnType\":4,\"columnLength\":30}"


10.136.0.58

curl -X POST \
  http://10.136.0.58:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:kingclean" \
  -H "x-zc-major-domain-id:5473" \
  -H "x-zc-sub-domain:r506" \
  -H "x-zc-sub-domain-id:5737" \
  -d "{\"columnName\":\"deviceSnCode\",\"columnType\":4,\"columnLength\":30}"

curl -X POST \
  http://10.136.0.58:5003/zc-bind/v1/addDeviceProperty \
 -H  "Content-Type:application/x-zc-object" \
  -H "x-zc-developer-id:260" \
  -H "x-zc-major-domain:kingclean" \
  -H "x-zc-major-domain-id:5473" \
  -H "x-zc-sub-domain:r506" \
  -H "x-zc-sub-domain-id:5737" \
  -d "{\"columnName\":\"deviceType\",\"columnType\":4,\"columnLength\":10}"　
```
