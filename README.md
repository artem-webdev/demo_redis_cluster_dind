# Демка запуска тестов в редис кластере с библиотекой testcontainers

## как запустить  
```
1) docker pull redis:7.2.4-alpine это желательно но не обязательно если проблемы с пуллом этого образа ,
 то спульте не менее по версии и укажите ваш образ в ./backend/tests/redis_test_containers/base.go в константе imageRedis 
2) cd ./deployments
3) docker compose up --build в первый раз потом можно docker compose up
```
## особености , не стал заморачиваться писать скрипт , просто определил в entrypoint композ файла
```
   entrypoint: >
      tini -- sh -c "
      // готовый скрипт от образа docker:27-dind
      dockerd-entrypoint.sh &
      // это что бы успел запуститься демон докера и прокинуть сеть с сокетом, если упадет по этоцй причине увелечте sleep
      // или напишите скрипт на свой вкус )
      sleep 20;
      // запуск тестов без кеша 
      go test -v -count=1 -failfast ./tests/redis_test_containers/...
      "
 ```


