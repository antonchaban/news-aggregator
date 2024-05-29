# News Alligator

<img height="250" src="https://tse2.mm.bing.net/th/id/OIG1.6SalGnQ.s83FWdg9MdLg?pid=ImgGn" width="250"/>

### Engineer name: Anton Chaban

## Summary

This project involves developing News Aggregator application, which will collect, process, and provide access to news
articles from multiple sources. Offers interface for accessing diverse news content, supporting various filtering
options such as source, keyword, and date range.

## Motivation

The primary motivation for developing the News Aggregator is to provide users with a straightforward tool for accessing
news articles from multiple sources. This project aims to streamline the process of collecting and retrieving news
content, making it easier for users to stay informed. The expected outcome is a robust instrument that can fetch,
filter, and deliver news articles efficiently based on user preferences.

# API Design

## Use case diagram

The following diagram illustrates the primary use cases of the News Aggregator application:

<img height="600" src="https://i.imgur.com/JHF2CVb.png" width="800"/>

## Use cases

The News Aggregator application will support the following use cases:

| Use case name  | Article Creation                                                                                                       |
|----------------|------------------------------------------------------------------------------------------------------------------------|
| Use case ID    | UC-01                                                                                                                  |
| Goals          | Aggregator creates article                                                                                             |
| Actors         | Aggregator                                                                                                             |
| Trigger        | Aggregator receives parsed Article entity                                                                              |
| Pre-conditions | Database (in memory or other) initialized, service received required DB, article was parsed and Create func was called |
| Flow of Events | 1. Service receives Article entity to Create.                                                                          |
|                | 2. Service calls create function                                                                                       |
|                | 3. Creation performed in specified DB                                                                                  |
| Extension      | -                                                                                                                      |
| Post-Condition | Article created and added to specified DB                                                                              |

| Use case name  | View all articles                                    |
|----------------|------------------------------------------------------|
| Use case ID    | UC-02                                                |
| Goals          | User see all available articles                      |
| Actors         | User                                                 |
| Trigger        | User calls specified command to view all articles    |
| Pre-conditions | -                                                    |
| Flow of Events | 1. User calls specified command                      |
|                | 2. The request is sent to the appropriate service    |
|                | 3. Articles that are found are displayed to the user |
| Extension      | -                                                    |
| Post-Condition | Articles that are found are displayed to the user    |

| Use case name  | View articles by keyword                                   |
|----------------|------------------------------------------------------------|
| Use case ID    | UC-03                                                      |
| Goals          | User see all available articles by specified keyword       |
| Actors         | User                                                       |
| Trigger        | User calls specified command to search articles by keyword |
| Pre-conditions | -                                                          |
| Flow of Events | 1. User calls specified command                            |
|                | 2. The request is sent to the appropriate service          |
|                | 3. Articles that are found are displayed to the user       |
| Extension      | -                                                          |
| Post-Condition | Articles that are found are displayed to the user          |

| Use case name  | View articles by source                                                      |
|----------------|------------------------------------------------------------------------------|
| Use case ID    | UC-04                                                                        |
| Goals          | User see all available articles by specified source                          |
| Actors         | User                                                                         |
| Trigger        | User calls specified command to search articles by source                    |
| Pre-conditions | -                                                                            |
| Flow of Events | 1. User calls specified command                                              |
|                | 2. The request is sent to the appropriate service                            |
|                | 3. Articles that are found are displayed to the user                         |
| Extension      | Search by source should be available only from the pool of available sources |
| Post-Condition | Articles that are found are displayed to the user                            |

| Use case name  | View articles by date range                                                                      |
|----------------|--------------------------------------------------------------------------------------------------|
| Use case ID    | UC-05                                                                                            |
| Goals          | User see all available articles by specified date range                                          |
| Actors         | User                                                                                             |
| Trigger        | User calls specified command to search articles by date range                                    |
| Pre-conditions | -                                                                                                |
| Flow of Events | 1. User calls specified command                                                                  |
|                | 2. The request is sent to the appropriate service                                                |
|                | 3. Articles that are found are displayed to the user                                             |
| Extension      | Search by date range should be available only if received data has specified format (YYYY-MM-DD) |
| Post-Condition | Articles that are found are displayed to the user                                                |

## Used models

### Article

``` go
type Article struct {
Id          int
Title       string
Description string
Link        string
Source      string
PubDate     time.Time
}
```

## Used repositories

### In Memory Repository

```go
type ArticleInMemory struct {
Articles []model.Article
nextID   int
}
```

### Interface for accessing articles

```go
type Article interface {
GetAll() ([]model.Article, error)
GetById(id int) (model.Article, error)
Create(article model.Article) (model.Article, error)
Delete(id int) error
GetByKeyword(keyword string) ([]model.Article, error)
GetBySource(source string) ([]model.Article, error)
GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error)
}
```

## Used services

### Article Service

We can use different repositories implementations for the ArticleService. For example, we can use the InMemoryRepository
We just need to send the repository to the service. In the future, we can easily switch to another repository

```go
func NewArticleService(articleRepo repository.Article) *ArticleService {
return &ArticleService{articleRepoInMem: articleRepo}
}
```

## CLI Design

The News Aggregator application will be accessible via a Command Line Interface (CLI). The CLI will provide users with
various commands to interact with the application, such as fetching news articles, filtering articles based on source,
keyword, and date range, and displaying articles in a user-friendly format.

And then we just need to add required repositories to the service and run the application

### Commands

The following commands will be available in the CLI:

- `-help` - Displays a list of available commands and their descriptions.

Expected output:

```
  -date-end string
        Specify the end date to filter the news by (format: YYYY-MM-DD).
  -date-start string
        Specify the start date to filter the news by (format: YYYY-MM-DD).
  -help
        Show all available arguments and their descriptions.
  -keywords string
        Specify the keywords to filter the news by.
  -sources string
        Select the desired news sources to get the news from. Supported sources: abcnews, bbc, washingtontimes, nbc, usatoday

```

- `-date-end <string>` - Specify the end date to filter the news by (format: YYYY-MM-DD).

Expected output:

```
List of articles:

ID: int
Title: string
Date: time
Description: string
Link: string
Source: string
  
```

- `-date-start <string>` - Specify the start date to filter the news by (format: YYYY-MM-DD).

Expected output:

```
List of articles:

ID: int
Title: string
Date: time
Description: string
Link: string
Source: string
  
```

- `-keywords <string>` - Specify the keywords to filter the news by.

Expected output:

```
List of articles:

ID: int
Title: string
Date: time
Description: string
Link: string
Source: string
  
```

- `-source <string>` - Select the desired news sources to get the news from. Supported sources:
  `abcnews, bbc, washingtontimes, nbc, usatoday`

Expected output:

```
List of articles:

ID: int
Title: string
Date: time
Description: string
Link: string
Source: string
  
```

## Example Usage

The following examples demonstrate how users can interact with the News Aggregator application using the CLI:

Execute command:

```.\news-aggregator.exe -keywords="war" -sources="bbc" -date-start="2024-05-18" ```

Output:

```
ID: 29
Title: Israel war cabinet minister vows to quit if there is no post-war plan for Gaza
Date: 2024-05-18 23:22:26 +0000 UTC
Description: Recent weeks have seen an increasingly public rift over how Gaza should be governed after the war.
Link: https://www.bbc.com/news/articles/cekkz82gnzgo
Source: BBC News

ID: 50
Title: Singer Libianca on the pressure to take sides in Cameroon war
Date: 2024-05-19 01:04:14 +0000 UTC
Description: The star of viral hit People talks to the BBC about getting death threats for waving a Cameroonian flag.
Link: https://www.bbc.com/news/articles/c8vz35911r9o
Source: BBC News
```