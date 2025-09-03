# Operator

Operator's health check

## Description

Эта реализация предоставляет:

- HTTP endpoints для health checks (/health, /ready, /live);
- Активные проверки с заданным интервалом;
- Поддержку стандартных проверок Kubernetes;
- Возможность отправки статуса во внешние системы;
- Graceful shutdown.

Вы можете настроить порт и интервал проверок через environment variables:

- `HEALTH_CHECK_PORT`: по умолчанию 8080;
- `HEALTH_CHECK_INTERVAL`: по умолчанию 30 секунд;
- `OPERATOR_NAME`: имя вашего оператора.


