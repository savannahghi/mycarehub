# Content API specification

This file gives an overview of the Content API(s) required by the application to perfom content management activities.

Content is an educative or informative material thats meant to enlighten users on various aspects. In this application, the content is more oriented towards health.

The API(s) schemas are defined in `GraphQL`.

## EndPoint Definitions

Base URLs:
- https://mycarehub-testing.savannahghi.org/ide

- https://mycarehub-staging.savannahghi.org/ide

- https://mycarehub-prod.savannahghi.org/ide

### Schema Inputs
```
input ShareContentInput {
	UserID:    String!
	ContentID: Int! 
	Channel:   String! 
}
```

### Schema Types
```
type ContentItemCategory {
  id: Int!
  name: String!
  iconUrl: String!
}
```

```
type Meta {
  totalCount: Int!
}
```

```
type ContentMeta {
  contentType: String!
  contentDetailURL: String!
  contentHTMLURL: String!
  slug: String!
  showInMenus: Boolean
  seoTitle: String
  searchDescription: String
  firstPublishedAt: String!
  locale: String
}
```

```
type HeroImage {
  ID: Int!
  title: String!
}
```

```
type HeroImageRendition {
  url: String!
  width: Int!
  height: Int!
  alt: String!
}
```

```
type Document {
  ID: Int!
  Document: DocumentData!
  meta: DocumentMeta!
}
```

```
type DocumentMeta {
  type: String!
  documentDetailUrl: String!
  documentDownloadUrl: String!
}
```

```
type GalleryImage {
  ID: Int!
  image: ImageDetail!
}
```
```
type ImageDetail {
  ID: Int!
  title: String!
  meta: ImageMeta!
}
```

```
type ImageMeta {
  type: String!
  imageDetailUrl: String!
  imageDownloadUrl: String!
}
```

```
type FeaturedMedia {
  ID: Int!
  url: String!
  title: String!
  type: String!
  duration: Float
  width: Int
  height: Int
  thumbnail: String
}
```

```
type Author {
  ID: String!
}
```

```
type CategoryDetail {
  ID: Int!
  categoryName: String!
  categoryIcon: String!
}
```

```
type DocumentData {
  ID: Int!
  title: String!
  meta: DocumentMeta!
}
```

```
type Content {
  items: [ContentItem!]!
  meta: Meta!
}
```

```
type ContentItem {
  ID: Int!
  title: String!
  date: String!
  meta: ContentMeta!
  intro: String!
  authorName: String!
  itemType: String!
  timeEstimateSeconds: Int
  body: String!
  heroImage: HeroImage
  heroImageRendition: HeroImageRendition
  likeCount: Int!
  bookmarkCount: Int!
  viewCount: Int!
  tagNames: [String!]!
  shareCount: Int!
  documents: [Document]
  author: Author!
  categoryDetails: [CategoryDetail]
  featuredMedia: [FeaturedMedia]
  galleryImages: [GalleryImage]
}
```



## Query Definitions

### Mutations
```
extend type Mutation {
  shareContent(input: ShareContentInput!): Boolean!
  bookmarkContent(userID: String!, contentItemID: Int!): Boolean!
  UnBookmarkContent(userID: String!, contentItemID: Int!): Boolean!
  likeContent(userID: String!, contentID: Int!): Boolean!
  unlikeContent(userID: String!, contentID: Int!): Boolean!
  viewContent(userID: String!, contentID: Int!): Boolean!
}
```

### Queries
```
extend type Query {
  getContent(categoryID: Int, Limit: String!): Content!
  listContentCategories: [ContentItemCategory!]!
  getUserBookmarkedContent(userID: String!): Content
  checkIfUserHasLikedContent(userID: String!, contentID: Int!): Boolean!
  checkIfUserBookmarkedContent(userID: String!, contentID: Int!): Boolean!
}
```

### 1. Mutations
#### 1.1. Share Content
Share content allows a user to share content with another user.
```
mutation shareContent($input: ShareContentInput!) {
  shareContent(input: $input)
}
```
Variables:
```
{"input":{
  "UserID": "userID",
    "ContentID": 6,
    "Channel": "SMS"
  }
}
```

#### 1.2. Bookmark Content
Bookmark content allows a user to bookmark ceetain content.
```
mutation bookmarkContent($userID: String!, $contentItemID: Int!){
  bookmarkContent(userID: $userID,contentItemID: $contentItemID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentItemID": 7
}
```

#### 1.3. UnBookmark Content
UnBookmark content allows a user to unbookmark a content.
```
mutation unbookmarkContent{
  UnBookmarkContent( 
    contentItemID: 6
  )
}
```

#### 1.4. Like Content
Like content allows a user to like certain content.
```
mutation likeContent($userID: String!, $contentID: Int!){
  likeContent(userID: $userID, contentID: $contentItemID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentID": 7
}
```

#### 1.5. UnLike Content
UnLike content allows a user to unlike certain content.
```
mutation unlikeContent($userID: String!, $contentID: Int!){
  unlikeContent(userID: $userID, contentID: $contentItemID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentID": 7
}
```
#### 1.6. View Content
View content updates the content viewed by a certain user.
```
mutation viewContent($userID: String!, $contentID: Int!){
  unlikeContent(userID: $userID, contentID: $contentID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentID": 7
}
```

### 2. Queries
#### 2.1. Get Content
This API fetches all the content from the Content Managemement System(CMS)
```
query {
  getContent(Limit: "5") {
    meta{
      totalCount
    }
    items {
      ID
      title
      date
      intro
      authorName
      tagNames
      meta{
        contentType
        contentType
        contentHTMLURL
        slug
        showInMenus
        seoTitle
        searchDescription
        firstPublishedAt
        locale
      }
      itemType
      timeEstimateSeconds
      body
      heroImage{
        ID
        title
      }
      heroImageRendition{
        url
        width
        height
        alt
      }
      likeCount
      bookmarkCount
      viewCount
      shareCount
      author {
        ID
      }
      documents {
        ID
        Document {
          ID
          title
        }
        meta{
          type
          documentDetailUrl
          documentDownloadUrl
        }
      }
      categoryDetails{
        ID
        categoryName
        categoryIcon
      }
      featuredMedia{
        ID
        url
        title
        type
        width
        height
        thumbnail
        duration
      }
      galleryImages{
        ID
        image{
          ID
          title
          meta{
            imageDetailUrl
            imageDownloadUrl
          }
        }
      }
    }
  }
}
```

#### 2.2. List Content Categories
This API only lists the content categories
```
query  {
  listContentCategories{
    id
    name
    iconUrl
  }
}
```

#### 2.3. Get User Bookmarked Content
This API fetches any content that has been bookmarked by the user.
```
query getUserBookmarkedContent($userID: String!){
  getUserBookmarkedContent(userID: $userID){
    items {
      ID
      title
      date
      intro
      authorName
      tagNames
      meta{
        contentType
        contentType
        contentHTMLURL
        slug
        showInMenus
        seoTitle
        searchDescription
        firstPublishedAt
        locale
      }
      itemType
      timeEstimateSeconds
      body
      heroImage{
        ID
        title
      }
      heroImageRendition{
        url
        width
        height
        alt
      }
      likeCount
      bookmarkCount
      viewCount
      shareCount
      author {
        ID
      }
      documents {
        ID
        Document {
          ID
          title
        }
        meta{
          type
          documentDetailUrl
          documentDownloadUrl
        }
      }
      categoryDetails{
        ID
        categoryName
        categoryIcon
      }
      featuredMedia{
        ID
        url
        title
        type
        width
        height
        thumbnail
        duration
      }
      galleryImages{
        ID
        image{
          ID
          title
          meta{
            imageDetailUrl
            imageDownloadUrl
          }
        }
      }
    }
    meta{
      totalCount
    }
  }
}
```
Variables:
```
{
  "userID": "userID"
}
```

#### 2.4. Check if user has liked content
This API perfoms a check to ascertain whether the user has liked content or not.
```
query checkIfUserHasLikedContent($userID: String!, $contentID: Int!){
   checkIfUserHasLikedContent(userID: $userID, contentID: $contentID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentID": 10000
}
```

#### 2.5. Check if user has bookmarked content
This API perfoms a check to ascertain whether the user has bookmarked content or not.
```
query checkIfUserBookmarkedContent($userID: String!, $contentID: Int!){
   checkIfUserBookmarkedContent(userID: $userID, contentID: $contentID)
}
```
Variables:
```
{
  "userID": "userID",
  "contentID": 10000
}
```
