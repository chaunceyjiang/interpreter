# interpreter
一个简易的解释器实现

- 变量绑定
- 整数类型、浮点数、布尔
- 数值计算
- 内置函数
- 函数“一等公民”
- 闭包
- 字符串

example
------

```bash
>>> let x = 1
>>> x
1
>>> let y = 2
>>> y
2
>>> let z = x + y
>>> z
3
>>>3-2.8
0.2
>>>5.3 + 5.7 + 5 + 5 - 10
11
>>>
>>>let add = fn(a, b) { return a + b; };
>>>add(3,4)
7
>>>let twice = fn(f,n){return f(f(n));};    
>>>let addTwo = fn(x){return x+2;};
>>>twice(addTwo,2);
6
>>
```