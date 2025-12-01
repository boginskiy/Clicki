# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Оптимизация кода
Оптимизация кода осуществлялась встроенными инструментами профилирования. Оптимизации был подвергнут handler, осуществляющий запись данных непосредственно в память компьютера во время работы сервиса.

Оптимизация была осуществлена путем номинального увеличения емкости структуры "Map" с 20 до 1024 элементов. Данное решение эффективно за счет снижения количества аллоцированный памяти, так же уменьшилось время отклика сервиса в момент переполнения Map (hash таблицы) и перехеширование структуры.

#### Результат
File: main
Build ID: dad7a541698c177dbf8fd95f800033713bd1ceb6
Type: inuse_space
Time: 2025-11-23 12:18:07 MSK
Showing nodes accounting for -8874.62kB, 81.08% of 10944.89kB total
      flat  flat%   sum%        cum   cum%
-5820.70kB 53.18% 53.18% -5820.70kB 53.18%  github.com/boginskiy/Clicki/internal/repository.(*RepositoryMapURL).CreateRecord
-1536.14kB 14.04% 67.22% -1536.14kB 14.04%  github.com/boginskiy/Clicki/internal/model.NewURLTb (inline)
   -1026kB  9.37% 76.59%    -1026kB  9.37%  runtime.allocm
  532.26kB  4.86% 71.73%   532.26kB  4.86%  github.com/boginskiy/Clicki/internal/repository.NewRepositoryMapURL
 -512.02kB  4.68% 76.41%  -512.02kB  4.68%  github.com/boginskiy/Clicki/internal/preparation.(*Functions).TakeAllBodyFromReq
 -512.02kB  4.68% 81.08% -8892.89kB 81.25%  github.com/boginskiy/Clicki/internal/service.(*URLServ).CreateURL
 -512.01kB  4.68% 85.76%  -512.01kB  4.68%  github.com/boginskiy/Clicki/pkg.Scramble
     512kB  4.68% 81.08%      512kB  4.68%  sync.(*Pool).pinSlow