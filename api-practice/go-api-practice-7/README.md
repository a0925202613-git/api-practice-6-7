# Go API Practice 7 - 圖書館書籍與借閱

與 go-api-practice-5 相同架構，題材為**書籍 (books)** 與**借閱紀錄 (borrowals)**。  
Handler 內為練習用註解，需自行實作商業邏輯。

## 功能重點

- **書籍**：CRUD；取得列表時支援 `?available=true|false` 篩選可借閱/已借出。
- **借閱**：  
  - 借書 `POST /borrowals`：需檢查該書 `available = true`，否則回 400。  
  - 還書 `POST /borrowals/:id/return`：僅當該筆借閱的 `returned_at` 為 NULL 時可還書，並將書設為可借閱。

## 啟動

```bash
createdb practice7
cp .env.example .env
# 編輯 .env 的 DATABASE_URL

go mod tidy
go run .
```

預設 port: **8082**。

## 練習項目

依 `handlers/` 內註解實作：  
GetBooks（含 available 篩選）、GetBookByID、CreateBook、UpdateBook、DeleteBook、GetBorrowals、GetBorrowalByID、CreateBorrowal（檢查可借閱）、ReturnBorrowal（僅未歸還可還書）。
