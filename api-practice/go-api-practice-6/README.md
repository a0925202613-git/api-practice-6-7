# Go API Practice 6 - 餐廳菜單與訂單

與 go-api-practice-5 相同架構，題材為**菜單 (menus)** 與**訂單 (orders)**。  
Handler 內為練習用註解，需自行實作商業邏輯。

## 功能重點

- **菜單**：CRUD（取得全部/單筆、新增、更新、刪除）。
- **訂單**：  
  - 取得訂單列表，支援 `?status=pending|completed|cancelled` 篩選。  
  - 新增訂單時需檢查 `menu_id` 是否存在。  
  - **取消訂單** `PATCH /orders/:id/cancel`：僅當 `status = pending` 時可取消。

## 啟動

```bash
# 建立 DB（例：practice6）
createdb practice6

cp .env.example .env
# 編輯 .env 的 DATABASE_URL

go mod tidy
go run .
```

預設 port: **8081**（與 practice-5 錯開）。

## 練習項目

依 `handlers/` 內註解實作：  
GetMenus、GetMenuByID、CreateMenu、UpdateMenu、DeleteMenu、GetOrders（含 status 篩選）、GetOrderByID、CreateOrder（menu 存在檢查）、CancelOrder（僅 pending 可取消）。
