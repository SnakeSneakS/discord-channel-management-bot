usecasesレイヤで定義したポートに対する実装を提供する(実態を記述する) 

**controller**と**gateway**と**presenter** 

| frameworks & drivers -> interface adapters -> application business rules |
|-|
| device,web -> **controllers** -> usecase (input)   |
| db -> **gateway** -> usecase                |
| ui -> **presenters** -> usecase (output)            |

