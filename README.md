# Go SQL DB

[中文](README-zh.md)

"Go SQL DB" is a relational database that supports SQL queries for research purposes. The main goal is to show the basic principles and key design of a relational database to database enthusiasts. Therefore, to easily understand, a lot of tricks but not very rigorous design was adopted, and the amount of code was controlled at about 2000 lines.

## Features

1. Pure Golang implementation, does not rely on any third-party packages. Goconvey was only introduced in unit tests
1. Unit test coverage ≈ 73.5%

#### Storage Engine
1. Special Thanks to [Let's Build a Simple Database](https://cstack.github.io/db_tutorial/)
1. Data retrieval structure based on B-Tree
1. Disk persistence engine based on 4KB paging
1. Close to POD (Plain Old Data) speed serialization & deserialization

#### SQL Parser
1. Tokenizer is implemented based on text/scanner
1. Support simple SELECT, INSERT syntax
   1. SELECT supports WHERE of numeric type
   1. Support LIMIT, but not support ORDER BY temporarily
1. If you want to know how the SQL Parser that can be used in the production environment is implemented, please refer to the [SQL Parser](https://github.com/auxten/postgresql-parser) that I stripped from CockroachDB and supports the SQL-2011 standard

#### Execution Planner
1. [Select Implementation](planner/select.go) based on Volcano Model
1. HTTP-based query and insert interface

## Known Issues

1. No DDL is implemented for the time being, only a fixed Schema
    ```go
    struct Row {
        Id uint32
        Sex byte
        Age uint8
        Username [32]byte
        Email [128]byte
        Phone [64]byte
    }
    ```
2. For limited support for SQL syntax, see [Test Cases](parser/parser_test.go)
3. Tokenizer is based on a clever implementation of the Golang language itself, there will be problems with the support of special characters in some strings, which can be solved by quoting strings with `"`

## How to run
1. Run
   ```bash
   go run . test.db
   ```
2. INSERT
   
   Execute `INSERT INTO table (id, username, email) VALUES (10, auxten, "auxtenwpc_gmailcom")`

   BY accessing: http://localhost:8080/exec?q=INSERT%20INTO%20table%20(id,%20username,%20email)%20VALUES%20(10,%20auxten,%20%22auxtenwpc_gmailcom%22)

3. SELECT
   
   Query `SELECT * FROM table WHERE id > 3 LIMIT 10`

   BY accessing: http://localhost:8080/query?q=SELECT%20*%20FROM%20table%20WHERE%20id%20%3E%203%20LIMIT%2010

## Thanks

1. [SQL-2011 SQL Parser](https://github.com/auxten/postgresql-parser)
1. [Marshal/Unmarshal Code generation](https://github.com/andyleap/gencode/)
1. [Document-oriented, embedded SQL database: genji](https://github.com/genjidb/genji)
1. [CockroachDB](https://github.com/cockroachdb/cockroach)
1. [Let's Build a Simple Database](https://cstack.github.io/db_tutorial/)

