# Design: Add Bubbletea TUI Framework

## Context

nippo-cli は現在シンプルな `fmt.Print` による出力と `promptui` によるプロンプト
を使用している。Issue #39 ではより豊かなターミナル UI の要望があり、Bubbletea
フレームワークの導入を提案する。

### 現状のアーキテクチャ

```text
Controller → Presenter → View → fmt.Print/promptui
                ↓
         ConsolePresenter
         (Progress/Complete/Suspend)
```

### 制約

- Clean Architecture の層分離を維持する必要がある
- DI パターン（samber/do/v2）との互換性を保つ
- 既存の ViewModel/ViewProvider パターンを活かす

## Goals / Non-Goals

### Goals

- Bubbletea を導入してインタラクティブな TUI を実現
- スピナーによる処理中の視覚的フィードバック
- プログレスバーによる進捗率表示（将来的な進捗計算追加を見据えて）
- 改善されたテキスト入力体験
- 一貫したスタイリング（Lipgloss）

### Non-Goals

- フルスクリーン TUI アプリケーション化
- マウスサポート
- 複雑なダッシュボード UI

## Decisions

### Decision 1: Bubbletea + Bubbles + Lipgloss の採用

**理由**: Charm エコシステムは Go TUI のデファクトスタンダードであり、MVU
アーキテクチャは Clean Architecture と親和性が高い。`bubbles` は一般的な UI
コンポーネント（spinner, textinput）を提供し、`lipgloss` は一貫したスタイリング
を可能にする。

**代替案**:

- `tview`: より重量級で、本プロジェクトには過剰
- `termui`: メンテナンスが活発でない
- `promptui` の継続使用: リッチ UI 機能が不足

### Decision 2: View 層での Bubbletea 統合

**理由**: Clean Architecture を維持するため、Bubbletea の詳細は View 層に
カプセル化する。Presenter は引き続き ViewModel を介して View と通信し、
Bubbletea の `tea.Program` 実行は View 内部で完結させる。

```text
Presenter → ViewModel → ViewProvider → Bubbletea Program
                                            ↓
                                     Update/View cycle
                                            ↓
                                     Channel feedback
```

### Decision 3: インラインモードの使用

**理由**: `tea.WithAltScreen()` によるフルスクリーンモードは本プロジェクトの
スコープ外。インラインモードでスピナーや入力を表示し、完了後は通常の出力に戻る。

### Decision 5: 終了時のUI出力保持

**理由**: Bubbletea はデフォルトで終了時に `View()` の最終出力を表示する。
現在の実装では `done` フラグが `true` の時に空文字を返しているため、UI がクリア
されてしまう。ユーザーは操作の経過を確認したいため、最終状態を保持する。

**実装方式**:

1. `View()` で `done` 状態でも最終メッセージを返す
2. 成功時: チェックマーク付きの完了メッセージを表示
3. 中断時 (Ctrl-C): 中断メッセージを表示
4. `done` と `interrupted` フラグで状態を管理し、`View()` 内で直接表示を生成
   （処理中と終了時で同じロジックを共有し、一貫性を保つ）

```go
func (m SpinnerModel) View() string {
    if m.done {
        if m.interrupted {
            return fmt.Sprintf("%s %s %s", WarningStyle.Render("✗"), m.message, DimStyle.Render("(interrupted)"))
        }
        return fmt.Sprintf("%s %s", SuccessStyle.Render("✓"), m.message)
    }
    return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}
```

### Decision 4: チャネルと Bubbletea Msg の統合

**理由**: 現在の実装は `chan<- interface{}` で View から Presenter に入力を返す。
Bubbletea は `tea.Cmd`/`tea.Msg` パターンを使用する。これらを統合するため：

```text
Presenter                     View (Bubbletea)
    |                              |
    |-- ViewModel (Input chan) --> |
    |                              | tea.Program.Run()
    |                              | User input → Update()
    |                              | tea.Quit
    | <-- chan <- result --------- |
    |                              |
```

- View は `tea.Program` を起動し、完了時に結果をチャネルに送信
- Presenter は従来通りチャネルから結果を受け取る
- Bubbletea の内部状態管理（Model）は View 層に閉じ込める

## Risks / Trade-offs

### Risk 1: 依存関係の増加

**影響**: `promptui` (1 パッケージ) から bubbletea エコシステム (3 パッケージ)
への変更。
**軽減策**: いずれも活発にメンテナンスされており、広く採用されている。

### Risk 2: 学習コスト

**影響**: MVU パターンへの習熟が必要。
**軽減策**: シンプルなユースケース（spinner, textinput）から始める。

## Migration Plan

1. 新しい依存関係を追加
2. TUI 基盤コードを `view/tui/` に作成
3. init コマンドの View を Bubbletea に移行
4. 他のコマンドにスピナーを追加
5. promptui 依存を削除

ロールバック: git revert で元に戻せる。各ステップは独立したコミットとする。

## Resolved Questions

- **スピナースタイル**: Dot (`⣾⣽⣻⢿⡿⣟⣯⣷`) を使用
- **適用範囲**: 全コマンド (init, build, deploy, clean, update) にスピナーを追加
- **移行方式**: promptui を完全削除し、bubbletea に一括移行
- **プログレスバー**: 実装する（現時点ではフォールバックでスピナー表示、将来の
  進捗計算追加に備える）
