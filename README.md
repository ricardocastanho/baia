# Real Estate Search Platform

Welcome to the Real Estate Search Platform, a Go-based project designed to simplify the search for real estate properties using Neo4j and Artificial Intelligence (AI). This platform aims to centralize property information from various sources, offering users a powerful tool to find the best properties at the right price.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)

## Features

- **Centralized Property Data**: Collects and stores property data from multiple real estate websites.
- **Advanced Search Capabilities**: Allows users to search for properties based on various criteria such as location, price, and property type.
- **AI-Powered Recommendations**: Utilizes machine learning algorithms to provide users with personalized property recommendations.
- **Price Analysis**: Historical price data is analyzed to help users make informed decisions about the best time to buy.
- **Neo4j Integration**: Leverages the power of graph databases to model relationships between properties, locations, and users.

## Technologies Used

- **Go**: The primary programming language used to build the backend of the platform.
- **Neo4j**: A graph database used to store and query property data.
- **Artificial Intelligence**: Machine learning techniques applied to provide personalized recommendations and price analysis.
- **Docker**: Used for containerization to ensure consistent development and deployment environments.

## Installation

### Prerequisites

- Go 1.19+ installed on your machine.
- Docker (for running the Neo4j database).
- [Neo4j](https://neo4j.com/download/) installed and running.

### Steps

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/real-estate-search-platform.git
    cd real-estate-search-platform
    ```

2. Install Go dependencies:

    ```sh
    go mod tidy
    ```

3. Run the application:

    ```sh
    go run main.go
    ```
