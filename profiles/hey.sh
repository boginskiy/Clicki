#!/bin/bash

# # hey
# for i in {1..30}; do
#     text_value="https://practicum.yandex-$i.ru"
#     hey -n 1 -m POST -H "Content-Type: text/plain" -d "$text_value" http://localhost:8080
# done

# curl &
for i in {1..50000}; do
    text_value="https://practicum.yandex-$i.ru"
    curl -X POST -H "Content-Type: text/plain" -d "$text_value" http://localhost:8080
done