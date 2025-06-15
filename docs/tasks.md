# Nerine 開発タスク一覧

## プロジェクト進捗

- ✅ プロジェクト基盤セットアップ
- ✅ HTTPサーバー基本実装
- 🔄 現在進行中: ドキュメント整備
- ⏳ 次のフェーズ: API実装

## 優先度別タスク

### 🔴 高優先度 (必須機能)

#### 認証・セキュリティ
- [ ] API キー認証ミドルウェア実装
- [ ] 環境変数バリデーション
- [ ] CORS設定

#### ログ・監視
- [ ] 構造化ログ (zap) セットアップ
- [ ] リクエスト/レスポンスログ
- [ ] エラーハンドリング統一

#### Repository Layer
- [ ] ArticleRepository インターフェース定義
- [ ] CategoryRepository インターフェース定義
- [ ] microCMS SDK integration

#### UseCase Layer
- [ ] GetArticlesUsecase 実装
- [ ] GetArticleByIDUsecase 実装
- [ ] GetCategoriesUsecase 実装
- [ ] GetArticlesByCategoryUsecase 実装

#### Handler Layer
- [ ] 記事一覧 API (`GET /api/v1/articles`)
- [ ] 記事詳細 API (`GET /api/v1/articles/:id`)
- [ ] カテゴリ一覧 API (`GET /api/v1/categories`)
- [ ] カテゴリ別記事一覧 API (`GET /api/v1/categories/:slug/articles`)

### 🟡 中優先度 (機能拡張)

#### 機能追加
- [ ] 人気記事一覧 API (`GET /api/v1/articles/popular`)
- [ ] 最新記事一覧 API (`GET /api/v1/articles/latest`)
- [ ] ページネーション機能
- [ ] 検索機能

#### パフォーマンス
- [ ] キャッシュ機能 (Redis)
- [ ] レスポンス圧縮
- [ ] レート制限

### 🟢 低優先度 (最適化・運用)

#### テスト
- [ ] ユニットテスト (gomock)
- [ ] 統合テスト
- [ ] パフォーマンステスト
- [ ] セキュリティテスト

#### デプロイ・インフラ
- [ ] Lambda用エントリーポイント実装
- [ ] Docker化
- [ ] CI/CD パイプライン
- [ ] 監視・アラート設定

## 完了タスク

### ✅ Phase 1: 基盤構築
- [x] Go module 初期化
- [x] プロジェクト構造設計
- [x] Entity定義 (Article, Category)
- [x] Echo HTTPサーバー基本実装
- [x] ヘルスチェックエンドポイント
- [x] 環境変数設定 (PORT)

## 次のアクション

1. **API認証ミドルウェア実装** - セキュリティの基盤
2. **構造化ログ設定** - 運用監視の基盤
3. **Repository実装** - データアクセス層
4. **UseCase実装** - ビジネスロジック層
5. **Handler実装** - API エンドポイント

## 注意事項

- 各タスクは Clean Architecture の原則に従って実装
- DDD パターンを維持
- テスト可能な設計を心がける
- 将来のLambda対応を考慮した実装
