该工具为kubernetes插件，提供了一种配置语言sclang与一个pod调度配置的代码生成工具scgen，该工具基于golang开发。

# scgen

pod调度配置的代码生成工具scgen可以根据手工编写的配置语言sclang生成相应pod的亲和度的yaml配置文件，可作为koordinator的一个plugin使用

| 输入            |输出|
|---------------|-----|
| 手工编写的sclang语句 |
| 需要插入配置的yaml文件 |已插入的yaml文件|

# sclang

- 配置语言sclang为pod之间的配置关系的自然语言表示，覆盖k8s的亲和度关系与阿里的自定义调度规则
  该配置语言的语法为：

1. 每一行为一条语句，第一个单词为required或者preferred，表述配置关系是强制的还是只是倾向性；
   primaryPod与subPod分别为需要配置亲和度的pod与之有关系的pod；configRelationship表示pod之间的亲和度配置关系，目前支持了依赖与排斥两种关系
2. 一条required语句基本结构为： required: primaryPod configRelationship subPod1,subPod2...
3. 一条preferred语句基本结构为： preferred:weight primaryPod configRelationship subPod1,subPod2...
4. primaryPod与subPod通过label:value来指代唯一的pod

# example

针对app:appA的调度配置需求为：
app:appA依赖app:appB
app:appA排斥app:appC,app:appD
app:appA可以与app:appE放在同一台机器，优先级为80（优先级为1-100）
app:appA可以与app:appF,app:appG放在同一台机器，优先级为60（优先级为1-100）
app:appA尽量不与pp:appH放在同一台机器，优先级为40（优先级为1-100）

转换为sclang：
required: app:appA & app:appB
required: app:appA ^ app:appC,app:appD
preferred:80 app:appA & app:appE
preferred:60 app:appA & app:appF,app:appG
preferred:40 app:appA ^ app:apH

app:appA的初始配置appA-pod-affinity.yaml为：

```yaml
metadata:
  name: app_test
  label:
    app: appA
apiVersion: v1
kind: pod
spec:

```

生成器读取编写的sclang，生成插入affinity的appA-pod-affinity.yaml

```yaml
apiVersion: v1
kind: pod
metadata:
  name: app_test
  label:
    app: appA
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        labelSelector:
          - matchExpressions:
              - key: app
                operator: In
                values:
                  - appB
          - matchExpressions:
              - key: app
                operator: NotIn
                values:
                  - appC
              - key: app
                operator: NotIn
                values:
                  - appD
        topologyKey: ""
      preferredDuringSchedulingIgnoredDuringExecution:
        preference:
          - weight: 80
            podAffinityTerm:
              labelSelector:
                - matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - appE
          - weight: 60
            podAffinityTerm:
              labelSelector:
                - matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - appF
                    - key: app
                      operator: In
                      values:
                        - appG
          - weight: 40
            podAffinityTerm:
              labelSelector:
                - matchExpressions:
                    - key: app
                      operator: NotIn
                      values:
                        - apH

```

# 问题

1. 生成的yaml中topologyKey缺失，在设计中是一个缺陷
2. 还未与koordinator作适配，还未开发



