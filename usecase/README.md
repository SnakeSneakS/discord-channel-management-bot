entitiesレイヤに依存してビジネスロジックを実行する責務をおう。
portをもち、このportはadapterで実装を差し替え可能(interfaceのみ記述し実態は別)。  

依存: 
interactor -> ports (inputとoutput) 

flow example: 
ui (input) -> presenter -> interactor -> controller -> web (output)