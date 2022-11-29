# Affinity的约束形式化表达

该文章简述了针对kubernetes的Pod配置文件中的affinity约束的形式化方法。

# affinity约束规则介绍

在kubernetes对象中，Affinity用于亲和度约束。

若要使用 Pod 间亲和性，可以使用 Pod 规约中的`.affinity.podAffinity`字段。使用 Pod 间反亲和性，可以使用 Pod
规约中的`.affinity.podAntiAffinity`字段。

podAffinity中的**`topologyKey`**
为最主要的字段，标志着pod需要在该拓扑域中调度，而具体调度到哪个pod，Affinity无法直接指定，而是根据该拓扑域中的集群情况而调度。若要调度到指定的node中，需要使用`nodeSelector`
字段而不是`Affinity`。

亲和性与反亲和性两者内部的数据结构都相同，都有两种类型的亲和度：

- `requiredDuringSchedulingIgnoredDuringExecution`
- `preferredDuringSchedulingIgnoredDuringExecution`

### 数据结构介绍：

[]**requiredDuringSchedulingIgnoredDuringExecution / preferredDuringSchedulingIgnoredDuringExecution{**

**labelSelector{**

[]**matchLabels**:{key,value}

[]**matchExpressions**{

**key**: string

**operato**r: string

**values**: []string

}

**}**

**Namespaces** []string

**TopologyKey** string

**}**

### pod约束

pod间使用`matchExpressions`（`LabelSelectorRequirement`
）作具体的约束。key为label的值，其中operator为label中基于集合的标签选择算符，有`in`、`notin`和`exists`三种

### 并列

其中有多个并列的slice关系，因此需要理清楚当其多个并列时是or还是and关系。

- 多个`requiredDuringSchedulingIgnoredDuringExecution`并列时：`and`，即需要都满足
- 多个`matchLabels`并列时：**`and`**，即需要都满足
- 多个**`matchExpressions`**并列时：**`and`**，即需要都满足
- 在一个matchExpressions中**`values`**有多个值时：**`or`**，即满足一个就可以了

### 案例

在某拓扑域中，有个name=app-a的pod需要调度，它有如下约束：

- `app-a`不能和标签为`app=app-b`的pod放在一起。（亦或者说，一个节点处于 Pod 所在的同一可用区且至少一个 Pod 具有`app=app-b`
  标签，则该 Pod 不应被调度到该节点上）（疑问，若app-a已经部署到某node上了，那么若app-b没有相关的亲和度配置，app-b可以成果部署在该node上吗
  ）
- app-a若要运行，该node中必须要有标签为`app=app-c`的pod。（亦或者说，仅当节点和至少一个已运行且有`app=app-c`的标签的 Pod
  处于同一区域时，才可以将该 Pod 调度到节点上。）

# 容错检查

## 案例

kubernetes对象的配置可以表现为yaml配置文件的形式。下面为一个yaml的案例，现在将其转换成CNF。

下面是一个yaml配置文件的一部分内容，主要为affinity部分。在其中有如下约束：

- app-a依赖于app=B或app=C
- app-a不能和app=C或app=D放在一起
- app-a依赖于infra=I1

可以看到实际上该配置是有冲突的。app-a依赖于app=C，同时又与app=C互斥，因此需要将其中的冲突检测出来。

```yaml
metadata:
  name: app-a
  labels:
    app: A
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
              - key: app
                operator: In
                values:
                  - B,C
            
              - key: app
                operator: NotIn
                values:
                  - C,D
              - key: infra
                operator: In
                values:
                  - I1
```

## yaml转换成SMT

smt-lib是smt solver的输入语言

根据上述的配置，将约束转换成SMT的形式。

$app=A → ((app=B \vee app=C) \wedge \neg (app=C \vee app=D) \wedge (infra=I1))$

将其用smtlib的格式表示出来,其中各个变量都是bool形式：

```scala
（declare-fun app=A () Bool）

（declare-fun app=B () Bool）

（declare-fun app=C () Bool）

（declare-fun app=D () Bool）

（declare-fun infra=I1 () Bool）

（assert（⇒  app=A (and（or app=B app=C）(and  (not (or app=B app=C)) infra=I1)）)

（check-sat）
```

## 规则

因此，kubernetes中affinity约束形式化的规则为：

1. 指定pod使用label
2. in转换成$→\neg$
3. notin转换成$→\neg$s
4. 多个macthExpressions使用and连接
5. label中多个value使用or连接