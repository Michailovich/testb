# Тестовое задание для проекта "Бани"

## Запуск сервисов

Соберите и запустите сервисы:

```bash
docker-compose up --build
```

### Создание записи

```
mutation {
  createMain(input: {title: "Test", sub_id: 1, sub_obj: "test"}) {
    id
    title
    created_at
  }
}
```

### Обновление записи:
```
mutation {
  updateMain(id: 1, input: {title: "Updated", deleted_at: ""}) {
    id
    title
    deleted_at
  }
}
```

### Удаление записи
```
mutation {
  deleteMain(id: 1)
}
```