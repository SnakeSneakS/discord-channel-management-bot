# discord channel management bot
discordのchannelを好きに作成したり参加・不参加したりできるbot 

## 問題点
- discord側でいじるとDBがそれに対応しない
    - 一致するように対応する、or 完全にdiscordから取得する(そもそもDBを使わない)、で基本はいい気がする. チャンネルがpublicかprivateかはdiscordに置けない独自仕様なのでDBに置くしかない気がするけど... 
- コマンド入力、UX的に渋いよね(かと言ってチャンネルの数だけリアクションつけるとかも渋いし...うーん。)

