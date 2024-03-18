# Guide to Using the `Version` Struct in Your Go Project

The `Version` struct is a crucial component in our Go project. It leverages the GORM (Go Object-Relational Mapper) library for database operations and the `ksuid` library for generating unique identifiers. This guide will help you understand how to use the `Version` struct effectively in your project.

## Adding the `Version` Type to Your Entity

Firstly, you need to add the `Version` type to your entity. This field will be automatically set before creating a new record, or you can set it manually to prevent reflection.

```go
type TestEntity struct {
    ...
	Version concurrency.Version
}
```

## Creating a New Record

When creating a new record, the `Version` field will be automatically set.

```go
e := TestEntity{
	ID:   1,
	Name: "1",
}
err := DB.Create(&e).Error
```

## Updating a Record

When updating a record, the `Version` field will be updated automatically. If the `Version` field in the database does not match the `Version` field in the entity, the update will fail. This is known as optimistic concurrency control.

```go
affected := DB.Model(&ec).Update("name", "3").RowsAffected
```

The generated SQL for the update operation will look like this:

```sql
UPDATE "test_entities" SET "name"='33',"version"='2dfAIbPnnndKFEAjxbpqYX2lCeX' WHERE "test_entities"."version" = '2dfAIb22lmBQ8grpFNvALvUX1Cq' AND "id" = 3
```

## Conclusion

The `Version` struct is a powerful tool for handling nullable string fields in a database and for implementing optimistic concurrency control. By understanding and utilizing its features, you can greatly enhance the functionality and reliability of your Go project.
```