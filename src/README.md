curl -X POST http://192.168.30.130:808/api/processParameter/report -d '{"xx":"xxxx"}'

curl -X POST http://127.0.0.1:808/api/processParameter/report -d '{"xx":"xxxx"}'




curl -X POST http://127.0.0.1:9888/po/deliver_goods -d '{"xx":"xxxx"}'




curl -X GET http://172.18.100.85/pdf
ab -n 1000000 -k -r -c 100 http://127.0.0.1/pdf
./weighttp -n 500000 -c 100 -t 2 -k 127.0.0.1/pdf

./weighttp -n 500000 -c 100 -t 2 -k 127.0.0.1:9888/redis