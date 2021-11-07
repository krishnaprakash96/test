--  Database
CREATE DATABASE ArticleDB;
USE ArticleDB;

-- Users table
CREATE TABLE
IF NOT EXISTS Articles(
    Title VARCHAR(32) NOT NULL,
    ArticleDescription VARCHAR(255) NOT NULL,
    Author VARCHAR(32) NOT NULL,
    Email VARCHAR(32) NOT NULL,
    CreatedTime DATETIME,
    PRIMARY KEY(Title));