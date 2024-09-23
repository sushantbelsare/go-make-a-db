
# SimpleDB

SimpleDB is a basic implementation of a relational database in Go, accessible via a command-line interface (CLI). It allows you to create tables, insert, select, update, and delete records using simple commands.

## Features

- Create and drop tables
- Insert records into tables
- Select records with optional conditions
- Update records based on conditions
- Delete records based on conditions
- Interactive CLI mode

## Getting Started

### Prerequisites

- Go 1.16 or later installed on your system

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/simple-db.git
   cd simple-db
   ```
2. Build the project:

   ```bash
   go build -o simpledb ./cmd
   ```

3. Run the application:

   ```bash
   ./simpledb
   ```

## Usage

Once the application is running, you can use the following commands in the CLI:

- **Create a table**: 
  ```
  create <table_name> <column1> <column2> ...
  ```
  Example:
  ```
  create users id name email
  ```

- **Drop a table**:
  ```
  drop <table_name>
  ```
  Example:
  ```
  drop users
  ```

- **List all tables**:
  ```
  list
  ```

- **Insert a record**:
  ```
  insert <table_name> <value1> <value2> ...
  ```
  Example:
  ```
  insert users 1 Alice alice@example.com
  ```

- **Select records**:
  ```
  select <table_name> [<column>=<value>]
  ```
  Example:
  ```
  select users name=Alice
  ```

- **Update records**:
  ```
  update <table_name> <col>=<val> <cond_col>=<cond_val>
  ```
  Example:
  ```
  update users email=alice@newdomain.com name=Alice
  ```

- **Delete records**:
  ```
  delete <table_name> <column>=<value>
  ```
  Example:
  ```
  delete users name=Alice
  ```

- **Show help**:
  ```
  help
  ```

- **Exit the program**:
  ```
  exit
  ```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project is inspired by learning exercises in building simple databases and CLI applications using Go.
