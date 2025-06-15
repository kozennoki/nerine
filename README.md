# Nerine - Blog BFF API

Next.js + microCMS ブログシステム用のBFF（Backend for Frontend）API

## 概要

microCMSをデータソースとして、Next.jsフロントエンドに最適化されたAPIを提供するGoアプリケーション

## アーキテクチャ

- **DDD + クリーンアーキテクチャ**
- **HTTPフレームワーク**: Echo
- **ログ**: zap（構造化ログ）
- **設定**: 環境変数（os.Getenv）
- **デプロイ**: AWS Lambda（将来対応）

## API エンドポイント

```
GET /api/v1/articles?page=1&limit=10          # 記事一覧（ページネーション）
GET /api/v1/articles/:id                      # 記事詳細
GET /api/v1/articles/popular?limit=5          # 人気記事一覧
GET /api/v1/articles/latest?limit=5           # 最新記事一覧
GET /api/v1/categories/:slug/articles?page=1  # カテゴリ別記事一覧
GET /api/v1/categories                        # カテゴリ一覧
```

### 認証

APIキーベース認証（Header: `X-API-Key`）

## microCMS APIスキーマ

### ブログ (endpoint: blog)
リスト形式のコンテンツタイプ

| フィールドID | 表示名 | 種類 | 必須 |
|-------------|--------|------|------|
| title | タイトル | テキストフィールド | true |
| category | カテゴリー | コンテンツ参照 - カテゴリー | true |
| body | 本文 | リッチエディタ | true |
| description | 概要 | テキストフィールド | true |
| image | 画像 | 画像 | true |

### カテゴリー (endpoint: categories)
リスト形式のコンテンツタイプ

| フィールドID | 表示名 | 種類 | 必須 |
|-------------|--------|------|------|
| name | カテゴリ名 | テキストフィールド | true |

## プロジェクト構造

```
/
├── cmd/
│   └── server/          # ローカル開発用サーバー
├── internal/
│   ├── domain/          # ビジネスルール
│   │   ├── entity/      # エンティティ
│   │   ├── service/     # ドメインサービス
│   │   └── repository/  # リポジトリインターフェース
│   ├── usecase/         # アプリケーションのユースケース
│   ├── infrastructure/  # 外部依存実装
│   │   ├── microcms/    # microCMS SDK wrapper
│   │   └── logger/      # zap logger
│   └── interfaces/      # コントローラー・プレゼンター
│       ├── handlers/    # Echo ハンドラー
│       ├── middleware/  # 認証ミドルウェア
│       └── presenter/   # レスポンス変換
└── pkg/                 # 外部から使用可能なパッケージ
```

## 環境変数

```bash
MICROCMS_API_KEY=your_microcms_api_key
MICROCMS_SERVICE_ID=your_microcms_service_id
NERINE_API_KEY=your_nerine_api_key
PORT=8080
```

## 開発開始

```bash
# Go モジュール初期化
go mod init github.com/your-username/nerine

# 依存関係インストール
go mod tidy

# サーバー起動
go run cmd/server/main.go
```

## 依存関係

- Echo: HTTPフレームワーク
- microCMS Go SDK: CMS API クライアント
- zap: 構造化ログ
