<script setup lang="ts">
const markdown = ref(`
# 一级标题

## 二级标题

### 三级

#### 四级

##### 五级

###### 六级

**加粗**

*倾斜*

***斜体加粗***

~~删除线~~

M^上标^

M~下标~

==高亮==

123123

\`行内代码\`

代码块
\`\`\`java
public class App {
    public static void main(String[] args) {
        System.out.println("Hello World!!");
    }
}
\`\`\`

\`\`\`c
#include <stdio.h>
int main()
{
    printf("Hello world!!");
}
\`\`\`

\`\`\`go
package main

import "fmt"

func main() {
    fmt.println("Hello world")
}
\`\`\`

+ 无序列表
+ ()
+ fd
+ ef

> y
>
> > 1
> >
> > > 1

---

| 123    | 123123  |
| ------ | ------- |
| 123    | 123     |
| 123123 | 13242   |
| 123123 | 123123  |
| 12312  | 1234142 |

$ a \\times b = c 公式$

1. 有序列表

<!-- 注释 -->

![demo](https://gitee.com/Immortal112/image/raw/master/background/11.webp)

[百度](https://baidu.com/)超链接 [名字]+(链接)

\`\`\`mermaid
graph LR
A-->b
C-->A
\`\`\`
`);

const src = ref("https://gitee.com/Immortal112/image/raw/master/background/11.webp");
const title = ref("Markdown 渲染示例");
</script>

<template>
  <ContentPanel :spacer="false">
    <article v-viewer class="article">
      <ArticleCover :src="src" :title="title" />
      <main class="content">
        <MarkdownRenderer :markdown></MarkdownRenderer>
      </main>
    </article>
  </ContentPanel>
</template>

<style lang="less" scoped>
.article {
  width: 100%;
}

.content {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  animation: article-show 1s ease-in-out;
}

@keyframes article-show {
  0% {
    opacity: 0;
    transform: translateY(32px);
  }
  100% {
    opacity: 1;
  }
}
</style>
