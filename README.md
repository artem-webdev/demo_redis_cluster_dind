# Демка запуска тестов в редис кластере с библиотекой testcontainers

## как запустить  
```
cd ./deployments && docker compose up --build 
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


