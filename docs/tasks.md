# Nerine 開発タスク一覧

## プロジェクト進捗

- ✅ プロジェクト基盤セットアップ
- ✅ HTTPサーバー基本実装
- ✅ 基本API実装完了
- ✅ gomock テスト基盤構築完了
- ✅ 基本ユニットテスト実装完了
- 🔄 現在進行中: テスト拡張・機能追加
- ⏳ 次のフェーズ: デプロイ準備・パフォーマンス最適化

## 優先度別タスク

### 🔴 高優先度 (必須機能)

#### テスト実装 (gomock)
- [x] gomock セットアップ (go.mod追加)
- [x] mock生成用タスクをTaskfile.ymlに追加
- [x] ArticleRepository mock生成
- [x] CategoryRepository mock生成
- [x] GetArticlesUsecase テスト実装
- [x] GetArticleByIDUsecase テスト実装  
- [x] GetCategoriesUsecase テスト実装
- [x] ArticleHandler テスト実装
- [x] CategoryHandler テスト実装
- [ ] UseCase interface mock生成 (Handler層テスト用)
- [ ] Handler層テスト拡張・エラーケース強化
- [ ] テストカバレッジ目標: 80%以上

#### 認証・セキュリティ
- [x] API キー認証ミドルウェア実装
- [x] 環境変数バリデーション
- [x] godotenv パッケージ導入（.env ファイル自動読み込み）
- [ ] CORS設定

#### ログ・監視
- [x] 構造化ログ (zap) セットアップ
- [ ] リクエスト/レスポンスログ
- [ ] エラーハンドリング統一

#### Repository Layer
- [x] ArticleRepository インターフェース定義
- [x] CategoryRepository インターフェース定義
- [x] microCMS SDK integration

#### UseCase Layer
- [x] GetArticlesUsecase 実装
- [x] GetArticleByIDUsecase 実装
- [x] GetCategoriesUsecase 実装
- [ ] GetArticlesByCategoryUsecase 実装

#### Handler Layer
- [x] 記事一覧 API (`GET /api/v1/articles`)
- [x] 記事詳細 API (`GET /api/v1/articles/:id`)
- [x] カテゴリ一覧 API (`GET /api/v1/categories`)
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

#### 高度なテスト
- [ ] 統合テスト (実際のmicroCMS API使用)
- [ ] パフォーマンステスト
- [ ] セキュリティテスト
- [ ] E2Eテスト

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

### ✅ Phase 2: API・テスト基盤構築
- [x] microCMS SDK integration
- [x] Repository層実装 (Article, Category)
- [x] UseCase層実装 (GetArticles, GetArticleByID, GetCategories)
- [x] Handler層実装 (基本API)
- [x] gomock環境セットアップ
- [x] Repository mock生成
- [x] UseCase層ユニットテスト実装
- [x] Handler層基本テスト実装
- [x] 構造化ログ (zap) セットアップ
- [x] API キー認証ミドルウェア実装

## 次のアクション

### 🎯 現在の最優先: テスト拡張・機能追加

1. **UseCaseインターフェースmock生成** - Handler層テスト強化
   - UseCase interfaceのmock生成
   - Handler層でのテスト拡張
   
2. **テストカバレッジ改善** - 品質向上
   - エラーケースの網羅的テスト
   - カバレッジ80%以上達成
   
3. **新機能実装** - API機能拡張
   - カテゴリ別記事一覧API実装
   - 追加エンドポイントのテスト実装

### 📋 実装順序
1. ✅ gomock環境構築 (完了)
2. ✅ Repository mock生成 (完了)
3. ✅ UseCase層テスト実装 (完了)
4. ✅ Handler層基本テスト (完了)
5. 🔄 UseCase interface mock生成
6. 🔄 Handler層テスト拡張
7. ⏳ テストカバレッジ測定・改善

## 注意事項

- 各タスクは Clean Architecture の原則に従って実装
- DDD パターンを維持
- テスト可能な設計を心がける
- 将来のLambda対応を考慮した実装
