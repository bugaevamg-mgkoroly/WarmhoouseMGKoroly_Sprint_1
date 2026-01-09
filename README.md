# Project_template
«Тёплый дом» — это небольшая компания, которая организует удалённое управление отоплением в доме. 

# Задание 1. Анализ и планирование

<aside>

Архитектура приложения представляет из себя монолит на Go с СУБД Postgres. Всё синхронно. Никаких асинхронных вызовов, микросервисов и реактивного взаимодействия в системе нет. Всё управление идёт от сервера к датчику. Данные о температуре также получаются через запрос от сервера к датчику.

</aside

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

Пользователи могут:

- включать и выключать отопление удалённо через приложение;
- устанавливать целевую температуру в помещении (в градусах Цельсия);
- выбирать предустановленные режимы работы («Комфорт», «Эконом», «Антизаморозка»);
- настраивать расписание работы системы по дням недели и временным интервалам;
- просматривать текущий статус системы (включено/выключено, установленный режим);
- получать уведомления об изменении статуса работы отопления.

Система поддерживает:

- синхронное управление одним или несколькими контурами отопления;
- контроль допустимого диапазона температур (минимальное/максимальное значение);
- автоматическое переключение между режимами по расписанию;
- принудительную остановку работы системы в аварийных ситуациях;
- передачу команд от сервера к исполнительным устройствам (котлам, сервоприводам);
- проверку доступности оборудования перед отправкой команд.

Ограничения:

- управление возможно только через централизованный сервер (нет прямого взаимодействия пользователь‑устройство);
- все операции выполняются синхронно — пользователь ждёт ответа системы;
- нельзя создавать пользовательские режимы работы (только выбор из предустановленных);
- отсутствует возможность группового управления несколькими объектами (один аккаунт — один дом).

Взаимодействие с другими компонентами:

- получает актуальные данные о температуре из модуля мониторинга;
- передаёт статусы выполнения команд в модуль уведомлений;
- запрашивает информацию об оборудовании из модуля учёта клиентов и устройств.

**Мониторинг температуры:**

Пользователи могут:

- просматривать текущую температуру в помещении в реальном времени;
- видеть статус подключения и работоспособности датчика (онлайн/офлайн);
- устанавливать пороги температуры для автоматических уведомлений;
- отслеживать температуру по отдельным зонам/комнатам (если установлено несколько датчиков);
- получать мгновенные уведомления при выходе температуры за заданные пределы.

Система поддерживает:

- регулярный опрос датчиков с заданным интервалом (например, каждые 5 мин);
- отображение данных в числовом и графическом виде (линейный график);
- индикацию аномалий (резкие скачки, длительные отклонения);
- мультизонный мониторинг (несколько датчиков в одном объекте);
- кэширование последних показаний для быстрого отображения;
- проверку достоверности данных (фильтрация шумов, контроль диапазонов).

Ограничения:

- данные доступны только за период работы подключённой системы;
- нет прогноза температуры или аналитических выводов;
- обновление показаний происходит с фиксированной частотой (не в режиме 100 % real‑time).

**Получение исторических данных о температуре**

Пользователи могут:

- запрашивать архивные данные за выбранный период (день, неделя, месяц);
- просматривать графики температуры с возможностью масштабирования по времени;
- экспортировать данные в формате CSV или PDF для анализа;
- сравнивать температурные профили за разные периоды;
- выделять фрагменты графика для детального изучения.

Система поддерживает:

- хранение исторических показаний с временной меткой;
- агрегацию данных (средние, минимальные, максимальные значения за период);
- фильтрацию по датчикам и зонам;
- пагинацию при выводе больших объёмов данных;
- индексацию по времени для быстрого поиска;
- ограничение глубины архива (например, до 1 года).

Ограничения:

- экспорт доступен только за периоды не более 30 дней подряд;
- детализация данных зависит от частоты опроса датчиков (например, только по 5‑минутным интервалам);
- нет возможности добавлять комментарии или метки к историческим данным.

**Организация сервисного обслуживания (выезды специалистов для подключения)**

Пользователи могут:

- подавать заявку на подключение системы через приложение;
- выбирать удобную дату и время выезда специалиста;
- отслеживать статус заявки (принята, назначена, выполнена);
- получать уведомления о приближении времени визита;
- оставлять отзыв о работе специалиста после завершения подключения;
- просматривать историю выполненных обслуживаний.

Система поддерживает:

- формирование заявки с автоматическим заполнением данных объекта (адрес, модель оборудования);
- проверку доступности слотов для выездов;
- назначение специалиста на заявку;
- отправку уведомлений пользователю и мастеру;
- фиксацию факта выполнения работ с фотоотчётом (опционально);
- учёт времени реагирования и длительности обслуживания;
- интеграцию с календарём для планирования графиков.

Ограничения:

- самостоятельное подключение оборудования недоступно (только через специалиста);
- изменение даты выезда возможно не позднее чем за 24 часа;
- нет онлайн‑трансляции процесса подключения;
- история обслуживаний хранится не более 2 лет.

### 2. Анализ архитектуры монолитного приложения

Монолитное приложение на Go с PostgreSQL в качестве СУБД.
Текущие особенности решения:
- Нынешнее приложение компании позволяет только управлять отоплением в доме и проверять температуру.
- Каждая установка сопровождается выездом специалиста по подключению системы отопления в доме к текущей версии системы.
- Архитектура приложения представляет из себя монолит на Go с СУБД Postgres. Всё синхронно. Никаких асинхронных вызовов, микросервисов и реактивного взаимодействия в системе нет. Всё управление идёт от сервера к датчику. Данные о температуре также получаются через запрос от сервера к датчику.
- Самостоятельно подключить свой датчик к системе пользователь не может.
  
### 3. Определение доменов и границы контекстов

Компании Теплый дом можно выделить следующие домены и поддомены:

Домен: управление отоплением
поддомен регулировка температуры по комнатам
	контекст: синхронные запросы к устройствам
	контекст: обработка ответов
поддомен расписание работы отопления (по дням/времени)
	контекст: планирование расписания
поддомен аварийное отключение
	контекст: отправка push‑уведомлений и email при авария
	
Домен: мониторинг температуры
поддомен отображение текущей температуры
	контекст: синхронные запросы к устройствам
поддомен графики изменений за период
	контекст: отображение собранных данных с датчиков
поддомен уведомления о критических значениях
	контекст: отправка push‑уведомлений
	
### **4. Проблемы монолитного решения**

- Высокий риск ошибок. Изменения в одной части приложения могут непредсказуемо влиять на другие части. Из-за этого возрастает вероятность возникновения ошибок и компании приходится тратить дополнительные ресурсы на тестирование.
- Длительные циклы разработки и развёртывания. При каждом изменении приходится тестировать всё приложение целиком. Это замедляет выпуск новых функций.
- Трудно управлять командой. В больших командах работа над монолитом часто приводит к конфликтам и задержкам. Изменения, которые вносит одна команда, влияют на работу других команд.
- Трудно масштабировать отдельные компоненты системы. Например, часть системы, которая отвечает за обработку заказов, испытывает высокую нагрузку. С монолитной архитектурой не получится масштабировать только эту часть — придётся масштабировать приложение целиком.

### 5. Визуализация контекста системы — диаграмма С4

[Контекст системы Теплый дом](https://disk.yandex.ru/d/IzFAsAebSbTg8A)

# Задание 2. Проектирование микросервисной архитектуры

**Диаграмма контейнеров (Containers)**

[Диаграмма контейнеров](https://disk.yandex.ru/d/2n8djUZydVa7hg)

@startuml C4_Container
!include <C4/C4_Container>


Container(api_gateway, "API Gateway", "Go, Kong/Traefik", "Единая точка входа, аутентификация (JWT), маршрутизация")
Container(control_svc, "Service‑Control", "Go", "Управление отоплением: команды, расписания, проверка устройств")
Container(monitoring_svc, "Service‑Monitoring", "Go", "Опрос датчиков, фильтрация шумов, кэширование")
Container(history_svc, "Service‑History", "Go, PostgreSQL", "Хранение и агрегация исторических данных")
Container(servicedesk_svc, "Service‑ServiceDesk", "Go, PostgreSQL", "Заявки на выезд, планирование, фиксация работ")
Container(notification_svc, "Service‑Notification", "Go", "Отправка push/email, шаблоны, повторные попытки")
ContainerDb(postgres_main, "PostgreSQL (основная БД)", "SQL", "Таблицы: devices, schedules, service_requests и др.")
Container(redis_cache, "Redis", "Кэш текущих показаний температуры")
Container(event_bus, "Event Bus", "Kafka/RabbitMQ", "Асинхронная коммуникация между сервисами")
Container(config_svc, "Config Service", "Consul", "Централизованная конфигурация (интервалы, пороги)")


System_Boundary(c1, "Экосистема «Тёплый дом»") {
  Rel(api_gateway, control_svc, "HTTP/REST, gRPC")
  Rel(api_gateway, monitoring_svc, "HTTP/REST")
  Rel(api_gateway, history_svc, "HTTP/REST")
  Rel(api_gateway, servicedesk_svc, "HTTP/REST")
  Rel(api_gateway, notification_svc, "HTTP/REST")


  Rel(control_svc, monitoring_svc, "gRPC: запрос текущей температуры")
  Rel(control_svc, event_bus, "Pub: command.executed, alarm")
  Rel(monitoring_svc, redis_cache, "Чтение/запись кэша")
  Rel(monitoring_svc, event_bus, "Pub: temperature.updated")
  Rel(event_bus, history_svc, "Sub: temperature.updated → запись в архив")
  Rel(event_bus, notification_svc, "Sub: alarm → отправка уведомления")
  Rel(servicedesk_svc, event_bus, "Pub: service.request.created")
  Rel(event_bus, notification_svc, "Sub: service.request.created → уведомление пользователя и мастера")
  Rel(control_svc, postgres_main, "SQL: чтение/запись устройств и расписаний")
  Rel(servicedesk_svc, postgres_main, "SQL: чтение/запись заявок и специалистов")
  Rel(history_svc, postgres_main, "SQL: запись исторических данных")
  Rel(config_svc, control_svc, "Получение конфигураций")
  Rel(config_svc, monitoring_svc, "Получение конфигураций")
}

@enduml

**Диаграмма компонентов (Components)**

[Диаграмма компонентов](https://disk.yandex.ru/d/4XTR4YuZnn8LAg)

@startuml
!include <C4/C4_Component.puml>

Component(temp_controller, "TemperatureController", "Go", "Обработка команд установки температуры")
Component(schedule_manager, "ScheduleManager", "Go", "Управление расписаниями работы")
Component(safety_checker, "SafetyChecker", "Go", "Контроль аварийных ситуаций (перегрев, замерзание)")
Component(device_commander, "DeviceCommander", "Go", "Формирование и отправка команд устройствам")
Component(validator, "Validator", "Go", "Проверка диапазонов температуры, прав доступа")

System_Boundary(c1, "Service‑Control") {
  Rel(temp_controller, validator, "Вызов: проверка допустимости температуры")
  Rel(temp_controller, device_commander, "Передача команды на устройство")
  Rel(schedule_manager, temp_controller, "Установка целевой температуры по расписанию")
  Rel(safety_checker, temp_controller, "Сигнал: аварийное отключение")
  Rel(device_commander, external_sensors, "gRPC/HTTP: отправка команды")
}

@enduml

**Диаграмма кода (Code)**

[Диаграмма кода](https://disk.yandex.ru/d/73QUEChvoVaMTQ)

@startuml C4_Code_TempController
class TemperatureController {
  + SetTargetTemp(deviceID: string, temp: float64) error
  + GetCurrentTemp(deviceID: string) (float64, error)
  - validateTempRange(temp: float64) bool
  - sendCommandToDevice(deviceID: string, cmd: string) error
}

# Задание 3. Разработка ER-диаграммы

[ER-модель](https://disk.yandex.ru/d/80D4UPzbPV-khw)

@startuml ER_Model_WarmHome
!include <C4/C4_Container.puml>

' Определяем сущности (таблицы)
entity "Device" as device {
  *id : UUID
  name : varchar(100)
  type : varchar(50)
  location : varchar(200)
  is_active : boolean
  last_seen : timestamp
}

entity "Schedule" as schedule {
  *id : UUID
  device_id : UUID
  start_time : timestamp
  end_time : timestamp
  target_temp : float
  repeat_pattern : varchar(50)
  is_enabled : boolean
}

entity "TemperatureReading" as temp_reading {
  *id : UUID
  device_id : UUID
  temperature : float
  timestamp : timestamp
  is_valid : boolean
}

entity "ServiceRequest" as service_request {
  *id : UUID
  device_id : UUID
  user_id : UUID
  title : varchar(200)
  description : text
  status : varchar(20) << (P, #FFAAAA) PENDING, IN_PROGRESS, COMPLETED, CANCELLED >>
  priority : int
  created_at : timestamp
  assigned_to : UUID
  completed_at : timestamp
}

entity "User" as user {
  *id : UUID
  username : varchar(50)
  email : varchar(100)
  role : varchar(20) << (A, #AAAAFF) ADMIN, USER, TECHNICIAN >>
  phone : varchar(20)
}

entity "Notification" as notification {
  *id : UUID
  user_id : UUID
  type : varchar(30) << (E, #AAFFAA) PUSH, EMAIL, SMS >>
  subject : varchar(200)
  body : text
  sent_at : timestamp
  success : boolean
}

entity "DeviceStatus" as device_status {
  *device_id : UUID
  is_overheated : boolean
  is_frozen : boolean
  alarm_triggered : boolean
  last_check : timestamp
}

' Определяем связи (relationships)
device ||--o{ schedule : "имеет расписания"
device ||--o{ temp_reading : "генерирует показания"
device ||--o{ device_status : "имеет статус"
device ||--o{ service_request : "связано с заявками"

user ||--o{ service_request : "создаёт заявки"
user ||--o{ notification : "получает уведомления"

service_request }o--|| user : "назначается технику"

' Индексы и внешние ключи (необязательные пояснения)
note as FK_NOTE
  Внешние ключи:
  - schedule.device_id → device.id
  - temp_reading.device_id → device.id
  - service_request.device_id → device.id
  - service_request.user_id → user.id
  - notification.user_id → user.id
  - device_status.device_id → device.id
end note
FK_NOTE .right. device

@enduml

# Задание 4. Создание и документирование API

### 1. Тип API

Для взаимодействия микросервисов в системе «Тёплый дом» оптимально комбинировать два типа API:
- HTTP/REST — для синхронных запросов «запрос‑ответ».
- Асинхронные сообщения через Event Bus (Kafka/RabbitMQ) — для событийно‑ориентированного взаимодействия.

1. HTTP/REST (синхронный обмен)
Где применяется:
- Взаимодействие клиент‑приложение ↔ API Gateway.
- Запросы от API Gateway к сервисным микросервисам (Service‑Control, Service‑Monitoring и др.).
- Прямые вызовы между сервисами, где нужен мгновенный ответ (например, Service‑Control → Service‑Monitoring: «дай текущую температуру»).

Почему именно REST:
- Простота и универсальность. HTTP‑методы (GET, POST, PUT, DELETE) интуитивно понятны, поддерживаются всеми языками и фреймворками.
- Кеширование и идемпотентность. Методы GET и PUT позволяют кешировать ответы и гарантировать идемпотентность операций.
- Широкая экосистема. Инструменты для мониторинга (Prometheus), трассировки (Jaeger), документации (OpenAPI/Swagger) работают «из коробки».
- Поддержка API Gateway. Kong/Traefik легко маршрутизируют и защищают REST‑эндпоинты (JWT, rate limiting).

Примеры вызовов:

GET /api/v1/devices/{id}/temperature → от Service‑Control к Service‑Monitoring.

POST /api/v1/service-requests → от клиента через API Gateway к Service‑ServiceDesk.

2. Асинхронные сообщения (Event Bus)
Где применяется:
- События типа «что‑то произошло» без ожидания ответа:
temperature.updated → запись в архив (Service‑History).
alarm → отправка уведомления (Service‑Notification).
service.request.created → оповещение пользователя и техника.

- Распределённые транзакции (например, обновление статуса устройства + создание заявки).
- Отложенные задачи (повторная отправка уведомлений, фоновая агрегация данных).

Почему Kafka/RabbitMQ:
- Отказоустойчивость. Сообщения сохраняются в очереди до обработки, даже если сервис временно недоступен.
- Масштабируемость. Несколько экземпляров сервиса могут потреблять сообщения параллельно (например, 3 инстанса Service‑Notification обрабатывают уведомления).
- Развязка сервисов. Производители и потребители не зависят друг от друга: можно добавлять новые сервисы‑подписчики без изменения кода отправителей.
- Гарантия доставки. Поддержка стратегий at‑least‑once или exactly‑once (в Kafka).
- Логирование событий. История сообщений позволяет анализировать цепочки событий (например, «как развивалась аварийная ситуация»).

Примеры событий:

{"event": "temperature.updated", "device_id": "abc123", "temp": 25.5, "timestamp": "2025-01-01T12:00:00Z"}

{"event": "alarm", "device_id": "abc123", "type": "overheat", "severity": "critical"}

Почему не только REST или только Event Bus?
- Только REST → риск «каскадных сбоев» (если один сервис тормозит, блокируются все зависимые), жёсткая связность, перегрузка сети частыми опросами.
- Только Event Bus → сложность отладки (нет явных запросов/ответов), избыточность для простых операций (например, «получить текущую температуру»).

Комбинация решает эти проблемы:
- REST — для быстрых, предсказуемых операций.
- Event Bus — для фоновых, критичных к надёжности задач.

Дополнительные технологии (опционально)
gRPC — для высокопроизводительных внутренних вызовов между сервисами (например, Service‑Control ↔ Service‑Monitoring), где важна скорость и низкая задержка.

WebSockets — для push‑уведомлений клиенту (например, мгновенное оповещение об аварии).

Итог
Гибридная модель API:
- REST — основной протокол для клиентских запросов и синхронных межсервисных вызовов.
- Event Bus (Kafka/RabbitMQ) — для асинхронных событий и распределённых процессов.

Это обеспечивает:
- баланс между скоростью, надёжностью и масштабируемостью;
- гибкость при добавлении новых сервисов;
- устойчивость к временным сбоям.

### 2. Документация API

Ниже приведены спецификации API для микросервисов системы «Тёплый дом» в формате OpenAPI 3.0 (Swagger) для REST‑конечных точек и AsyncAPI 2.4.0 для событийной шины.
1. REST API (OpenAPI 3.0)
Сервис: Service‑Control
Файл: service-control-openapi.yaml

yaml
openapi: 3.0.3
info:
  title: Service‑Control API
  version: 1.0.0
  description: Управление отоплением: команды, расписания, проверка устройств.
servers:
  - url: https://api.warmhome.example.com/control

paths:
  /devices/{deviceId}/target-temperature:
    put:
      summary: Установить целевую температуру
      operationId: setTargetTemperature
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                temperature:
                  type: number
                  format: float
                  minimum: -30
                  maximum: 60
      responses:
        '200':
          description: Температура установлена
        '400':
          description: Некорректное значение температуры
        '404':
          description: Устройство не найдено

  /schedules:
    post:
      summary: Создать расписание работы устройства
      operationId: createSchedule
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [deviceId, startTime, endTime, targetTemp]
              properties:
                deviceId:
                  type: string
                startTime:
                  type: string
                  format: date-time
                endTime:
                  type: string
                  format: date-time
                targetTemp:
                  type: number
                  format: float
      responses:
        '201':
          description: Расписание создано
        '400':
          description: Ошибки валидации
Сервис: Service‑Monitoring
Файл: service-monitoring-openapi.yaml

yaml
openapi: 3.0.3
info:
  title: Service‑Monitoring API
  version: 1.0.0
  description: Опрос датчиков, фильтрация шумов, кэширование.
servers:
  - url: https://api.warmhome.example.com/monitoring

paths:
  /devices/{deviceId}/current-temperature:
    get:
      summary: Получить текущую температуру с датчика
      operationId: getCurrentTemperature
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Температура получена
          content:
            application/json:
              schema:
                type: object
                properties:
                  temperature:
                    type: number
                    format: float
                  timestamp:
                    type: string
                    format: date-time
        '404':
          description: Датчик не найден
Сервис: Service‑ServiceDesk
Файл: service-servicedesk-openapi.yaml

yaml
openapi: 3.0.3
info:
  title: Service‑ServiceDesk API
  version: 1.0.0
  description: Заявки на выезд, планирование, фиксация работ.
servers:
  - url: https://api.warmhome.example.com/servicedesk

paths:
  /service-requests:
    post:
      summary: Создать заявку на обслуживание
      operationId: createServiceRequest
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [deviceId, title, description]
              properties:
                deviceId:
                  type: string
                title:
                  type: string
                  maxLength: 200
                description:
                  type: string
                priority:
                  type: integer
                  minimum: 1
                  maximum: 5
      responses:
        '201':
          description: Заявка создана
        '400':
          description: Ошибки валидации

    get:
      summary: Получить список заявок
      operationId: listServiceRequests
      parameters:
        - name: status
          in: query
          schema:
            type: string
            enum: [PENDING, IN_PROGRESS, COMPLETED, CANCELLED]
      responses:
        '200':
          description: Список заявок
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: string
                    title:
                      type: string
                    status:
                      type: string
2. Event Bus API (AsyncAPI 2.4.0)
Файл: event-bus-asyncapi.yaml

yaml
asyncapi: 2.4.0
info:
  title: WarmHome Event Bus
  version: 1.0.0
  description: Асинхронные события системы «Тёплый дом».

servers:
  kafka:
    url: kafka://broker.warmhome.example.com:9092
    protocol: kafka

channels:
  temperature.updated:
    subscribe:
      summary: Событие обновления температуры
      message:
        payload:
          type: object
          properties:
            deviceId:
              type: string
            temperature:
              type: number
              format: float
            timestamp:
              type: string
              format: date-time

  alarm:
    subscribe:
      summary: Событие аварийной ситуации
      message:
        payload:
          type: object
          properties:
            deviceId:
              type: string
            type:
              type: string
              enum: [overheat, freeze, communication_loss]
            severity:
              type: string
              enum: [low, medium, high, critical]
            timestamp:
              type: string
              format: date-time

  service.request.created:
    subscribe:
      summary: Событие создания заявки на обслуживание
      message:
        payload:
          type: object
          properties:
            requestId:
              type: string
            deviceId:
              type: string
            title:
              type: string
            timestamp:
              type: string
              format: date-time
Как использовать
Для REST API (OpenAPI):

Загрузите YAML‑файлы в  для визуализации и тестирования.

Сгенерируйте клиентский код (на Go, Python, JS и др.) через .

Для Event Bus (AsyncAPI):

Используйте  для просмотра схемы событий.

Настройте консьюмеры/продюсеры в сервисах на основе описанных каналов (temperature.updated, alarm и др.).

Примечания
Версионирование: Все API имеют версию 1.0.0. В будущем добавляйте новые версии без нарушения совместимости.

Безопасность: Для REST API требуется JWT‑аутентификация (укажите в securitySchemes при необходимости).

Валидация: Схема OpenAPI включает ограничения (minimum, maximum, enum), которые можно использовать для автоматической валидации запросов.

Формат дат: Все временные метки в формате ISO 8601 (YYYY‑MM‑DDTHH:MM:SSZ).

# Задание 5. Работа с docker и docker-compose

Перейдите в apps.

Там находится приложение-монолит для работы с датчиками температуры. В README.md описано как запустить решение.

Вам нужно:

1) сделать простое приложение temperature-api на любом удобном для вас языке программирования, которое при запросе /temperature?location= будет отдавать рандомное значение температуры.

Locations - название комнаты, sensorId - идентификатор названия комнаты

```
	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
```

2) Приложение следует упаковать в Docker и добавить в docker-compose. Порт по умолчанию должен быть 8081

3) Кроме того для smart_home приложения требуется база данных - добавьте в docker-compose файл настройки для запуска postgres с указанием скрипта инициализации ./smart_home/init.sql

Для проверки можно использовать Postman коллекцию smarthome-api.postman_collection.json и вызвать:

- Create Sensor
- Get All Sensors

Должно при каждом вызове отображаться разное значение температуры

Ревьюер будет проверять точно так же.


# **Задание 6. Разработка MVP**

Необходимо создать новые микросервисы и обеспечить их интеграции с существующим монолитом для плавного перехода к микросервисной архитектуре. 

### **Что нужно сделать**

1. Создайте новые микросервисы для управления телеметрией и устройствами (с простейшей логикой), которые будут интегрированы с существующим монолитным приложением. Каждый микросервис на своем ООП языке.
2. Обеспечьте взаимодействие между микросервисами и монолитом (при желании с помощью брокера сообщений), чтобы постепенно перенести функциональность из монолита в микросервисы. 

В результате у вас должны быть созданы Dockerfiles и docker-compose для запуска микросервисов. 
