# Data models

The data model defines a number of related entities all used for the website.

```mermaid
erDiagram
  Item_Notes {
    int ID PK
    int Item_ID FK
    enum Item_Type
    int Note_ID FK
  }

  Notes {
    int Note_ID PK
    date Published_Date
    string Title
    string Subtitle
    text Body
  }

  Roles {
    int Role_ID PK
    date Start_Date
    date End_Date
    string Title
    string Subtitle
    int Company_ID FK
    string Company_Icon
    string Slug
    text Description
    int Skills
  }

  Companies {
    int Company_ID PK
    string Name
    text Description
    string Icon
  }

  Projects {
    int Project_ID PK
    date Start_Date
    date End_Date
    string Title
    text Description
  }

  Tagged_Items {
    int ID PK
    int Tag_ID FK
    int Item_ID FK
    enum Item_Type
  }

  Tags {
    int ID PK
    string Name
    string Slug
    string Icon
    string Theme
  }

  Item_Notes }|--o| Notes : "Note_ID"
  Roles }|--o| Companies : "Company_ID"
  Tagged_Items }o--o| Tags : "Tag_ID"
  Item_Notes }|--o| Projects : "Note_ID"
  Tagged_Items }o--o| Projects : "Item_ID"
  Tagged_Items }o--o| Notes : "Item_ID"
  Tagged_Items }o--o| Roles : "Item_ID"
  Item_Notes }|--o| Roles : "Note_ID"
```

# Schema
The SQL for the database is defined below:
```sql
CREATE TABLE IF NOT EXISTS companies (
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  description text,
  icon varchar(2083)
);

CREATE TABLE IF NOT EXISTS roles (
  id bigserial PRIMARY KEY,
  createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  deletedAt timestamp(0) with time zone,
  startDate timestamp(0) with time zone NOT NULL,
  endDate timestamp(0) with time zone,
  title text NOT NULL,
  subtitle text,
  companyId bigint NOT NULL,
  FOREIGN KEY (companyId) REFERENCES companies(id),
  slug varchar(255),
  description text,
  skills text[] NOT NULL
);

CREATE TABLE IF NOT EXISTS projects (
  id bigserial PRIMARY KEY,
  createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  deletedAt timestamp(0) with time zone,
  startDate timestamp(0) with time zone NOT NULL,
  endDate timestamp(0) with time zone,
  title text NOT NULL,
  description text
);

CREATE TABLE IF NOT EXISTS tags (
  id bigserial PRIMARY KEY,
  createdAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  deletedAt timestamp(0) with time zone,
  name text NOT NULL,
  slug varchar(255),
  icon varchar(1),
  theme varchar(255)
);

CREATE TABLE IF NOT EXISTS notes (
  id bigserial PRIMARY KEY,
  createdAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  deletedAt timestamp(0) with time zone,
  publishedAt timestamp(0) with time zone,
  title text NOT NULL,
  subtitle text,
  body text
);

CREATE TYPE item_type AS ENUM ('notes', 'roles', 'projects');

CREATE TABLE IF NOT EXISTS tagged_items (
  id bigserial PRIMARY KEY,
  tagId bigint NOT NULL,
  itemId bigint NOT NULL,
  itemType item_type NOT NULL,
  FOREIGN KEY (tagId) REFERENCES tags(id)
);

CREATE TABLE IF NOT EXISTS item_notes (
  id bigserial PRIMARY KEY,
  noteId bigint NOT NULL,
  itemId bigint NOT NULL,
  itemType item_type NOT NULL
);

```

## Notes table migration

To create the standalone table that powers the notes feature, apply the following SQL:

```sql
CREATE TABLE IF NOT EXISTS notes (
  id bigserial PRIMARY KEY,
  createdAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  deletedAt timestamp(0) with time zone,
  publishedAt timestamp(0) with time zone,
  title text NOT NULL,
  subtitle text,
  body text
);
```
