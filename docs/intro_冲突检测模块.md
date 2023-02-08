# 冲突检测模块

冲突检测模块将静态分析配置有冲突的语句，并报出冲突提示配置的逻辑冲突。
在整体的代码生成器中，该模块在配置静态优化之后，在yanl写入之前执行。

因此，该模块的输入与输出为：

| input | output                   | 
|-------|--------------------------|
| Pod对象 | 错误信息：1. 是否有冲突 2. 发生冲突的位置 |

## 应用场景

例如，在一个配置文件中 a in b，a notin b同时出现，需要报出这两条有冲突的错误。
具体来说，如下 pod.yaml:

```
.....
matchExpressions:
  - key: app
    operator: NotIn
    values:
      - appb
 ...... 
matchExpressions:
  - key: app
    operator: In
    values:
      - appb
```

该配置中，appb既In又NotIn，显然是逻辑会错误,因此需要报出其错误并定位到该两个语句。

## 工作原理

冲突检测的原理是将配置约束转化成一个可满足问题，使用求解器对该可满足问题进行求解。若该问题是不满足的，说明约束之间有冲突。
建模的变量为一个Pod对象或者配置文件的affinity中的具体约束。
因此，该sat问题可以表述为：
//TODO 建模的公式

如下图所示，冲突检测模块主要分为三步：

1. CNF转换
   读取Pod对象或者yaml文件，按照[Affinity的约束形式化表达](Affinity的约束形式化表达.md)的规则定义，按照

> A. requiredDuringSchedulingIgnoredDuringExecution/preferredDuringSchedulingIgnoredDuringExecution

> B. matchLabels/matchExpressions

> C. values

三级关系，转化成

```
A1 & A2 & A3 & ... 
( B11 & B12 & B13 & ...)
( C111 | C112 | C113 | ...)
```

因此,一个配置yaml中的约束关系可以表合取范式CNF的形式：
> (C111 | C112 | C113 ) & ( C121 | C122 ) & (C211 ) ...

其中`Cijk`用`label:value`的字符串标识其唯一性 ，在一个matchExpressions或者matchLabel中的value为CNF中的每一个析取字句。

通过以上的方法，将pod对象或者yaml建模成了一个可满足性问题。

2. Dimacs转换

> Dimacs是一种面向行的格式，由3种不同的基本类型的行组成的CNF表示语言。
> 三种类型的语句如下：
> 注释行：任何以“c”开头的行都是注释行。
> 摘要行：此行包含有关文件中问题的类型和大小的信息。摘要行以“p”开头，接着是问题的类型（在大多数情况下是“cnf”）、变量的数量和该问题中的子句的数量。
> 子句行：子句行由以空格分隔的数字组成，以0结尾。每个非零数字表示一个文字，负数是该变量的负文字，0是一行的终止符。

步骤1中的CNF的每一个Cijk都对应着字句行中的一个数字，因此用一个哈希表存储与转换其关系，对每个变量进行编码，生成Dimacs语句文件。

3. Solver求解
   使用现代sat/smt求解器求解其可满足性。
   当前使用minisat求解。调用minisat，可验证其满足性。
   当前的约束只需或与非的形式逻辑就可以表示，因此直接选用sat solver求解。

![冲突检测模块数据流图](/pic/冲突检测模块数据流图.jpg)

## 冲突定位
若一个文件有语句冲突，需要定位其文件中出现冲突的语句对的位置。
方法：先从现有的sat求解器中探究有无冲突定位的功能，如果没有，需要修改sat求解器的功能

## 容错机制

无

## 待完善点

1. 若之后的需求中有如“CpuSet=1“类似的约束，需要重新建模，使用谓词逻辑去表达，转换成SMT。使用SMT solver去求解。
2. 无冲突定位功能

## 参考链接
1.[Conflict-driven clause learning](https://en.wikipedia.org/wiki/Conflict-driven_clause_learning)
2.[tutorial: Conflict Driven Clause Learning University of Washington | CSE 442](https://cse442-17f.github.io/Conflict-Driven-Clause-Learning/)