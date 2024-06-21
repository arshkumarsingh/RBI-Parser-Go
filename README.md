# RBI Parser Go

RBI Parser Go is a Go-based project designed to parse data related to banks from the Reserve Bank of India (RBI). This tool can be used to extract and process bank information efficiently.

## Features

- Parses RBI data to extract bank information.
- Provides a JSON output of the parsed data.
- Can be easily integrated into other Go projects.

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/arshkumarsingh/RBI-Parser-Go.git
    cd RBI-Parser-Go
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

## Usage

1. Run the parser:

    ```sh
    go run cmd/main.go
    ```

2. The parsed bank data will be saved in a JSON file.

## Project Structure

- `cmd/`: Contains the main executable for running the parser.
- `download/`: Directory for any downloaded data files.
- `banks.json`: Sample JSON file containing bank data.
- `etags.json`: File for managing entity tags for HTTP caching.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries or issues, please contact the repository owner.
