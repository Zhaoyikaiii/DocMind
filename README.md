# DocMind

DocMind is an intelligent document analysis and question-answering system built with Golang and LangChain. It allows users to upload documents, process them, and perform intelligent queries to extract relevant information.

## Features

- **Document Upload**: Supports multiple document formats including PDF, Word, and Markdown.
- **Intelligent Analysis**: Utilizes AI to analyze and vectorize document content.
- **Question-Answering**: Provides accurate answers to user queries based on document content.
- **RESTful API**: Offers a robust API for integration with other applications.

## Tech Stack

- **Backend**: Golang
- **AI Integration**: LangChain
- **Database**: PostgreSQL for metadata, Milvus/Weaviate for vector storage
- **Framework**: Gin for HTTP server
- **Document Processing**: Apache Tika or native parsers

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker (for containerization)
- PostgreSQL
- Milvus/Weaviate (for vector storage)