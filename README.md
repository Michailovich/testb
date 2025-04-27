# Тестовое задание для проекта "Бани"

## Запуск сервисов

Соберите и запустите сервисы:

```bash
docker-compose up --build
```

### Создание записи

```
mutation {
  createMain(input: {
    title: "My Tool Main"
    tool: {
      title: "Hammer"
      description: "Steel claw hammer"
    }
  }) {
    id
    title
    subObj
    subId
    tool {
      id
      title
      description
      mainId
      createdAt
      updatedAt
      deletedAt
    }
    createdAt
    updatedAt
    deletedAt
  }
}
```

### Обновление записи:
```
mutation {
  updateMain(id: 1, input: {
    title: "Updated Main Title"
  }) {
    id
    title
    updatedAt
    deletedAt
  }
}
```

### Удаление записи(soft delete)
```
mutation {
  updateMain(id: 1, input: {
    deletedAt: "2024-07-01T12:00:00Z"
  }) {
    id
    title
    deletedAt
  }
}
```