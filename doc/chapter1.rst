目的
----

ソフトウェアルータ「Lagopus2」の基本機能実装を基盤とし、汎用サーバ上で動作し、IPsec機能やシステム冗長化機能などの高付加価値機能を有する、SDNベースの高性能・高機能ソフトウェアルータの実現と評価を行うため、Lagopus2データストア部の応用機能の試作を行う。
データストア部の応用機能の実現のため、データプレーン部のフレームワークの改造およびスイッチ機能を提供するための各モジュールの実装を行う。
本ドキュメントは、データプレーン部の新規フレームワークのアーキテクチャの説明および構造について説明する。また、各モジュールの構造についても説明する。

本開発物の特徴を以下に示す。

-  データプレーン部内でのシステム情報制御部とパケット処理部の分離
-  機能毎に独立したモジュール定義
-  外部エージェントからのシステム情報設定用APIの提供

なお、本開発では、前回の成果物であるLagopus2データストア部イテレーション1で作成されたソースコードをベースに実装を行う。
