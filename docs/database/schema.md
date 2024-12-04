# Database Schema Design

## Overview

This document describes the database schema design for the DocMind system. The system uses PostgreSQL as the primary database and implements features such as user management, document management, version control, and tagging system.

## Entity Relationship Diagram

```mermaid
erDiagram
    User ||--o{ Document : creates
    User ||--o{ DocumentVersion : creates
    User ||--o{ UserSettings : has
    Document ||--o{ DocumentVersion : has
    Document }|--|| User : owned_by
    Document }o--o{ Tag : has
    Document ||--o{ Document : has_parent

```

Table Structures
Users Table
Stores user information and authentication details.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(128) NOT NULL UNIQUE,
    full_name VARCHAR(128) NOT NULL,
    bio VARCHAR(256),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

User Settings Table
Stores user preferences and settings.

```sql
CREATE TABLE user_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    theme VARCHAR(20) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    notifications BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

Documents Table
Stores document content and metadata.

```sql
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    version INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'draft',
    creator_id INTEGER NOT NULL,
    parent_id INTEGER,
    path VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id),
    FOREIGN KEY (parent_id) REFERENCES documents(id)
);
```

Document Versions Table
Stores document version history.

```sql
CREATE TABLE document_versions (
    id SERIAL PRIMARY KEY,
    document_id INTEGER NOT NULL,
    version INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

```

Tags Table
Stores document tags.

```sql
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

Document Tags Table
Manages many-to-many relationships between documents and tags.

```sql
CREATE TABLE document_tags (
    document_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (document_id, tag_id),
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
```