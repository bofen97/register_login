#### 设计思路

-   Table struct
    -   负责和数据库的连接 CURD等

-   server  module
    -   Table struct 使用这个struct来和数据库做交互

-   serverData module
    -   负责http请求json数据的解析


