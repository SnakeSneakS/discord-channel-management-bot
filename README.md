# discord channel management bot
discordのchannelを好きに作成したり参加・不参加したりできるbot 

## 問題点
- discord側でいじるとDBがそれに対応しない
    - 一致するように対応する、or 完全にdiscordから取得する(そもそもDBを使わない)、で基本はいい気がする. チャンネルがpublicかprivateかはdiscordに置けない独自仕様なのでDBに置くしかない気がするけど... 
- コマンド入力、UX的に渋いよね(かと言ってチャンネルの数だけリアクションつけるとかも渋いし...うーん。)


## TODO: 
1. Update, Archive の実装
2. メンバー数や、自分の参加是非などの表示

## limitation: 
1. 1カテゴリに対して50までのチャンネルしか作れない (TODO: 複数カテゴリ用意する)
2. 